package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func main() {
	root, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	http.Handle("/", http.FileServer(http.Dir(root+"/assets")))
	http.HandleFunc("/spruce", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/spruce" {
			w.WriteHeader(404)
			w.Write([]byte("not found"))
			return
		}

		if r.Method != "POST" {
			w.WriteHeader(415)
			w.Write([]byte("method not supported"))
			return
		}

		all, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("failed to read POST body: %s", err)
			w.WriteHeader(400)
			w.Write([]byte("bad request"))
			return
		}

		var params struct {
			Prune []string
			YAML  []struct {
				Filename string
				Contents string
			}

			Debug bool
			Trace bool
		}
		err = json.Unmarshal(all, &params)
		if err != nil {
			log.Printf("failed to unmarshal JSON payload: %s", err)
			log.Printf("bad JSON payload was:\n%s", string(all))
			w.WriteHeader(400)
			w.Write([]byte("bad request"))
			return
		}

		dir, err := ioutil.TempDir("", "web")
		if err != nil {
			log.Printf("failed to create temporary working directory: %s", err)
			w.WriteHeader(500)
			w.Write([]byte("internal server error"))
			return
		}
		// defer rm() the dir

		var args []string
		if params.Debug {
			args = append(args, "--debug")
		}
		if params.Trace {
			args = append(args, "--trace")
		}
		args = append(args, "merge")

		if len(params.Prune) > 0 {
			args = append(args, "--prune")
			for _, f := range params.Prune {
				args = append(args, f)
			}
		}
		for _, y := range params.YAML {
			args = append(args, y.Filename)
			ioutil.WriteFile(
				fmt.Sprintf("%s/%s", dir, y.Filename),
				[]byte(y.Contents),
				os.FileMode(0400),
			)
		}

		cmd := exec.Command("spruce", args...)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Printf("failed to get stdout pipe from command: %s", err)
			w.WriteHeader(500)
			w.Write([]byte("internal server error"))
			return
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			log.Printf("failed to get stderr pipe from command: %s", err)
			w.WriteHeader(500)
			w.Write([]byte("internal server error"))
			return
		}

		err = os.Chdir(dir)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("internal server error"))
			return
		}

		log.Printf("running `spruce %v' in %s", args, dir)
		err = cmd.Start()
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("internal server error"))
			return
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
		var response struct {
			Arguments []string `json:"args"`
			Stdout    string   `json:"stdout"`
			Stderr    string   `json:"stderr"`
			Success   bool     `json:"success"`
		}
		response.Success = (err != nil)
		response.Arguments = args
		response.Stdout = <-output
		response.Stderr = <-errors

		m, err := json.Marshal(response)
		if err != nil {
			log.Printf("failed to marshal JSON response: %s", err)
			w.WriteHeader(500)
			w.Write([]byte("internal server error"))
			return
		}
		w.WriteHeader(200)
		w.Write(m)
	})
	panic(http.ListenAndServe(":8081", nil))
}
