package main

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/google/go-github/github"
)

func main() {
	wd, err := os.Getwd()
	log.Println(wd)
	if err != nil {
		panic(err)
	}
	user := flag.String("username", "eyedeekay", "username to generate pages for")
	flag.Parse()
	if _, err := os.Stat("config.json"); err != nil {
		jsonStruct := generate(*user)
		bytes, err := json.MarshalIndent(jsonStruct, "", "  ")
		if err != nil {
			panic(err)
		}
		if err := ioutil.WriteFile("config.json", bytes, 0644); err != nil {
			panic(err)
		}
	}

	if bytes, err := ioutil.ReadFile("config.json"); err != nil {
		panic(err)
	} else {
		reposList := make(map[string]string)
		json.Unmarshal(bytes, &reposList)
		for index, remote := range reposList {
			log.Println("git clone", remote, filepath.Join(wd, index))
		}
	}
}

// this is super freaking crude, don't do things like this.
func generate(gh_user string) map[string]string {
	client := github.NewClient(nil)
	jsonStruct := make(map[string]string)
	// list public repositories for org "github"
	length := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25}
	for i := range length {
		log.Println("integer", i)
		opt := &github.RepositoryListOptions{
			Type:      "public",
			Sort:      "updated",
			Direction: "asc",
			ListOptions: github.ListOptions{
				PerPage: 100,
				Page:    i,
			},
		}
		repos, _, err := client.Repositories.List(context.Background(), gh_user, opt)
		if err != nil {
			panic(err)
		}
		for _, repo := range repos {
			if repo.GetParent() == nil {
				if !*repo.Fork {
					if *repo.Name != gh_user+".github.io" {
						jsonStruct[*repo.Name] = repo.GetSSHURL()
					}
				}
			}
		}
	}
	return jsonStruct
}
