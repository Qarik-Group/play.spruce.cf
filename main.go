package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	Spruces []string
)

func init() {
	n := 5
	if s := os.Getenv("SPRUCES_BACK"); s != "" {
		v, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			log.Printf("unable to parse SPRUCE_BACK=%s: %s", s, err)
		} else {
			n = int(v)
		}
	}

	l, err := GetSpruces(n)
	if err != nil {
		panic(fmt.Sprintf("init failed: %s", err))
	}
	Spruces = l
}

func main() {
	go boom()

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

	http.HandleFunc("/meta", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/meta" {
			w.WriteHeader(404)
			w.Write([]byte("not found"))
			return
		}

		if r.Method != "GET" {
			w.WriteHeader(415)
			w.Write([]byte("method not supported"))
			return
		}
		meta := struct {
			Flavors []string `json:"flavors"`
		}{
			Flavors: Spruces,
		}

		res, err := json.Marshal(meta)
		if err != nil {
			log.Printf("failed to marshal JSON response: %s", err)
			w.WriteHeader(500)
			w.Write([]byte("internal server error"))
			return
		}
		w.WriteHeader(200)
		w.Write(res)
	})

	http.HandleFunc("/mem", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/mem" {
			w.WriteHeader(404)
			w.Write([]byte("not found"))
			return
		}

		switch r.Method {
		case "POST":
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

			key, err := Store(&m)
			if err != nil {
				log.Printf("failed to store Merge op in back-end store: %s", err)
				w.WriteHeader(500)
				w.Write([]byte("internal server error"))
				return
			}

			w.WriteHeader(200)
			w.Write([]byte(key))
			return

		case "GET":
			q := r.URL.Query()
			key, ok := q["k"]
			if !ok {
				w.WriteHeader(400)
				w.Write([]byte("bad request"))
				return
			}
			m, err := Retrieve(key[0])
			if err != nil {
				w.WriteHeader(404)
				w.Write([]byte("not found"))
				return
			}

			res, err := json.Marshal(m)
			if err != nil {
				log.Printf("failed to marshal JSON response: %s", err)
				w.WriteHeader(500)
				w.Write([]byte("internal server error"))
				return
			}
			w.WriteHeader(200)
			w.Write(res)

		default:
			w.WriteHeader(415)
			w.Write([]byte("method not supported"))
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("listening on :%s\n", port)
	panic(http.ListenAndServe(":"+port, nil))
}
