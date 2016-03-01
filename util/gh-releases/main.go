package main

import (
	"os"
	"fmt"
	"github.com/jhunt/play.spruce.cf/github"
)

func usage() {
	fmt.Fprintf(os.Stderr, "USAGE: gh-releases (list|download) [owner [repo]]\n")
	os.Exit(1)
}

func main() {
	owner, repo := "geofffranks", "spruce"
	if len(os.Args) < 1 || len(os.Args) > 4 {
		usage()
	}
	fmt.Printf("os.Args: %v\n", os.Args)
	switch os.Args[1] {
	case "list":
	case "download":
	default: usage()
	}
	if len(os.Args) > 2 {
		owner = os.Args[2]
	}
	if len(os.Args) > 3 {
		repo = os.Args[3]
	}

	r, err := github.Releases(owner, repo)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}
	versions := github.LatestFrom("1.0.2", r)

	switch os.Args[1] {
	case "list":
		for _, v := range versions {
			fmt.Printf("%s\n", v)
		}

	case "download":
		for _, v := range versions {
			fmt.Printf("downloading spruce v%s, to spruce-%s\n", v, v)
			f, err := os.Create(fmt.Sprintf("spruce-%s", v))
			if err != nil {
				fmt.Printf("error: %s\n", err)
				continue
			}
			github.Download(owner, repo, v, f)
			f.Close()
		}

	default:
		usage()
	}
}
