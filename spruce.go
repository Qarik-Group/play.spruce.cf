package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Result struct {
	Arguments []string `json:"args"`
	About     string   `json:"about"`
	Stdout    string   `json:"stdout"`
	Stderr    string   `json:"stderr"`
	Success   bool     `json:"success"`
}

func Spruce(where string, args ...string) (*Result, error) {
	r := &Result{Arguments: args}

	cmd := exec.Command("spruce", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("failed to get stdout pipe from command: %s", err)
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Printf("failed to get stderr pipe from command: %s", err)
		return nil, err
	}

	if where != "" {
		err = os.Chdir(where)
		if err != nil {
			log.Printf("failed to chdir to %s: %s", where, err)
			return nil, err
		}
	}
	where, _ = os.Getwd()
	log.Printf("running `spruce %v' in %s", args, where)
	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	output := make(chan string)
	errors := make(chan string)
	drain := func(rd io.Reader, out chan string) {
		var buf []string
		s := bufio.NewScanner(rd)
		for s.Scan() {
			buf = append(buf, s.Text()+"\n")
		}
		out <- strings.Join(buf, "")
		close(out)
	}
	go drain(stdout, output)
	go drain(stderr, errors)

	err = cmd.Wait()
	r.Success = (err != nil)
	r.Stdout = <-output
	r.Stderr = <-errors
	return r, nil
}

type Merge struct {
	Prune []string
	YAML  []struct {
		Filename string
		Contents string
	}

	Debug bool
	Trace bool
}

func SpruceMerge(m Merge) (*Result, error) {
	dir, err := ioutil.TempDir("", "web")
	if err != nil {
		log.Printf("failed to create temporary working directory: %s", err)
		return nil, err
	}
	// defer rm() the dir

	var args []string
	if m.Debug {
		args = append(args, "--debug")
	}
	if m.Trace {
		args = append(args, "--trace")
	}
	args = append(args, "merge")

	if len(m.Prune) > 0 {
		for _, f := range m.Prune {
			args = append(args, "--prune")
			args = append(args, f)
		}
	}
	for _, y := range m.YAML {
		args = append(args, y.Filename)
		ioutil.WriteFile(
			fmt.Sprintf("%s/%s", dir, y.Filename),
			[]byte(y.Contents),
			os.FileMode(0400),
		)
	}

	result, err := Spruce(dir, args...)
	if version, err := Spruce("", "-v"); err == nil {
		result.About = version.Stdout + version.Stderr
	} else {
		log.Printf("failed to determine spruce version information: %s", err)
	}
	return result, err
}
