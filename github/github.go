package github

import (
	"encoding/json"
	"fmt"
	"github.com/jhunt/play.spruce.cf/client"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"runtime"
	"strings"
)

var (
	VersionMatch *regexp.Regexp
)

func init() {
	VersionMatch = regexp.MustCompile(`v\d+\.\d+\.\d+`)
}

func Releases(owner, repo string) ([]string, error) {
	c := client.New()
	r, err := c.Get(fmt.Sprintf("https://api.github.com/repos/%s/%s/tags", owner, repo))
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, fmt.Errorf("API %s", r.Status)
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var data []struct{ Name string }
	err = json.Unmarshal(b, &data)
	if err != nil {
		return nil, err
	}

	versions := make([]string, 0)
	for _, tag := range data {
		if VersionMatch.MatchString(tag.Name) {
			versions = append(versions, tag.Name[1:])
		}
	}
	return versions, nil
}

func Latest(versions []string) []string {
	cut, n, last := make([]string, 0), 0, ""

	for i := range versions {
		ver := strings.Split(versions[i], ".")
		mm := strings.Join(ver[0:2], ".")

		if last != mm {
			last = mm
			n++
			if n > 1 {
				cut = append(cut, versions[i])
				continue
			}
		}
		if n <= 1 {
			cut = append(cut, versions[i])
		}
	}

	return cut
}

func LatestFrom(from string, versions []string) []string {
	vv := make([]string, 0)
	for _, v := range versions {
		vv = append(vv, v)
		if v == from {
			break
		}
	}
	return Latest(vv)
}

func Download(owner, repo, version string, out io.Writer) error {
	// use default http client always for downloading, as it sends to S3, and S3 doesn't like the Auth Headers.. plus this isn't rate-limited?
	r, err := http.Get(fmt.Sprintf("https://github.com/%s/%s/releases/download/v%s/spruce-%s-%s", owner, repo, version, runtime.GOOS, runtime.GOARCH))
	if err != nil {
		return err
	}
	if r.StatusCode != 200 {
		return fmt.Errorf("API %s", r.Status)
	}
	io.Copy(out, r.Body)
	return nil
}
