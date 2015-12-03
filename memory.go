package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type GistFile struct {
	Content string `json:"content"`
}
type Gist struct {
	ID          string              `json:"id,omitempty"`
	Description string              `json:"description"`
	Public      bool                `json:"public"`
	Files       map[string]GistFile `json:"files"`
}

func parseGist(r *http.Response) (*Gist, error) {
	raw, err := ioutil.ReadAll(r.Body)
	gist := Gist{}
	err = json.Unmarshal(raw, &gist)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response: %s", err)
	}

	return &gist, nil
}

func Store(m *Merge) (string, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	gist := Gist{
		Description: "a bit o' spruce, from http://play.spruce.cf",
		Public:      true,
		Files: map[string]GistFile{
			"spruce.json": GistFile{Content: string(data)},
		},
	}
	encoded, err := json.Marshal(gist)
	if err != nil {
		return "", err
	}
	body := strings.NewReader(string(encoded))
	resp, err := http.Post("https://api.github.com/gists", "application/json", body)
	if err != nil {
		return "", err
	}
	log.Printf("got a response: %v", resp)
	if resp.StatusCode != 201 {
		return "", fmt.Errorf("failed to create gist: API returned %s", resp.Status)
	}

	created, err := parseGist(resp)
	if err != nil {
		return "", err
	}
	return created.ID, nil
}

func Retrieve(key string) (*Merge, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/gists/%s", key))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("not found")
	}

	gist, err := parseGist(resp)
	if err != nil {
		return nil, err
	}

	f, ok := gist.Files["spruce.json"]
	if !ok {
		return nil, err
	}

	var m Merge
	err = json.Unmarshal([]byte(f.Content), &m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
