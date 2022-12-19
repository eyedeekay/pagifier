package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/github"
)

func main() {
	ex, err := os.Executable()
	if err != nil {
		log.Printf("Path error: %s", err)
	}
	wd := filepath.Dir(ex)
	user := flag.String("username", "eyedeekay", "username to generate pages for")
	flag.Parse()
	if _, err := os.Stat("config.json"); err != nil {
		jsonStruct := generate(*user)
		bytes, err := json.MarshalIndent(jsonStruct, "", "  ")
		if err != nil {
			log.Printf("Marshal error: %s", err)
		}
		if err := ioutil.WriteFile("config.json", bytes, 0644); err != nil {
			log.Printf("Write error: %s", err)
		}
	}

	if bytes, err := ioutil.ReadFile("config.json"); err != nil {
		log.Printf("Read error: %s", err)
	} else {
		reposList := make(map[string]string)
		json.Unmarshal(bytes, &reposList)
		for index, remote := range reposList {
			time.Sleep(time.Second)
			log.Println("git clone", remote, filepath.Join(wd, index))
			if _, err := os.Stat(filepath.Join(wd, index)); os.IsNotExist(err) {
				_, err := git.PlainClone(filepath.Join(wd, index), false, &git.CloneOptions{
					URL:      remote,
					Progress: os.Stdout,
				})
				if err != nil {
					log.Printf("Clone error: %s", err)
					continue
				}
			} else {
				r, err := git.PlainOpen(filepath.Join(wd, index))
				if err != nil {
					log.Printf("Open error: %s", err)
					continue
				}
				w, err := r.Worktree()
				if err != nil {
					log.Printf("Tree error: %s", err)
					continue
				}
				err = w.Pull(&git.PullOptions{RemoteName: "origin"})
				if err != nil {
					log.Printf("Pull error: %s", err)
					continue
				}
				ref, err := r.Head()
				if err != nil {
					log.Printf("Head error: %s", err)
					continue
				}
				commit, err := r.CommitObject(ref.Hash())
				if err != nil {
					log.Printf("Log error: %s", err)
					continue
				}
				fmt.Println(commit)
			}
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
			log.Printf(" error: %s", err)
			continue
		}
		if len(repos) == 0 {
			break
		}
		for _, repo := range repos {
			if repo.GetParent() == nil {
				if !*repo.Fork {
					if *repo.HasPages {
						log.Printf("repo %s has pages: %v", *repo.URL, *repo.HasPages)
						if *repo.Name != gh_user+".github.io" {
							jsonStruct[*repo.Name] = strings.Replace(repo.GetGitURL(), "git://", "https://", 1)
						}
					}
				}
			}
		}
	}
	return jsonStruct
}
