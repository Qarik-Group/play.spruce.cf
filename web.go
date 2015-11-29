package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

		var m Merge
		err = json.Unmarshal(all, &m)
		if err != nil {
			log.Printf("failed to unmarshal JSON payload: %s", err)
			log.Printf("bad JSON payload was:\n%s", string(all))
			w.WriteHeader(400)
			w.Write([]byte("bad request"))
			return
		}

		result, err := SpruceMerge(m)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("internal server error"))
			return
		}
		res, err := json.Marshal(result)
		if err != nil {
			log.Printf("failed to marshal JSON response: %s", err)
			w.WriteHeader(500)
			w.Write([]byte("internal server error"))
			return
		}
		w.WriteHeader(200)
		w.Write(res)
	})
	panic(http.ListenAndServe(":8081", nil))
}
