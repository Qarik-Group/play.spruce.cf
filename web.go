package main

import (
	"io/ioutil"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("assets")))
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
			w.WriteHeader(400)
			w.Write([]byte("bad request"))
			return
		}
		w.WriteHeader(200)
		w.Write([]byte(all))
	})
	panic(http.ListenAndServe(":8081", nil))
}
