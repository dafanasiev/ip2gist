package main

import (
	"context"
	"encoding/json"
	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Configuration struct {
	GetIpUrl string
	Gist struct {
		Id string
		AccessToken string
	}
}

func readConfig() Configuration{
	file, err := os.Open("conf.json")
	if err != nil || file == nil {
		log.Fatalf("cant read rv: %v", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	rv := Configuration{}
	err = decoder.Decode(&rv)
	if err != nil {
		log.Fatalf("invalid rv: %v", err)
	}
	return rv
}

func getMyIp(configuration Configuration) string {
	getIpResp, err := http.Get(configuration.GetIpUrl)
	if err != nil {
		log.Fatalf("Cant get IP addr from %s: %v", configuration.GetIpUrl, err)
	}
	defer getIpResp.Body.Close()

	if getIpResp.StatusCode > 299 {
		log.Fatalf("Cant get IP addr from %s, due to status code %v", configuration.GetIpUrl, getIpResp.StatusCode)
	}

	getIpRespBytes, err := ioutil.ReadAll(getIpResp.Body)
	if err != nil {
		log.Fatalf("Cant read response body from %s: %v", configuration.GetIpUrl, err)
	}

	rv := strings.TrimSpace(string(getIpRespBytes))
	log.Printf("My current public IP: %s", rv)

	return rv
}

func syncGistContent(configuration Configuration, myPublicIP string) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: configuration.Gist.AccessToken },
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	gist, _, err := client.Gists.Get(context.Background(), configuration.Gist.Id)
	if err != nil {
		log.Fatalf("Cant read gist with id %s: %v", configuration.Gist.Id, err)
	}

	updatedGistFiles := make(map[github.GistFilename]github.GistFile)
	for gistFileKey,gistFile := range gist.Files {
		cachedIp := strings.TrimSpace(*gistFile.Content)
		log.Printf("Cached IP: %s",cachedIp)
		if myPublicIP != cachedIp {
			log.Printf("Cached IP not equals Current Public IP - need update gist")
			gistFile.Content = &myPublicIP
			updatedGistFiles[gistFileKey] = gistFile
		}
	}

	if len(updatedGistFiles)>0 {
		log.Printf("Try to save new version of gist...")
		_,_, err := client.Gists.Edit(context.Background(), configuration.Gist.Id,  &github.Gist{Files: updatedGistFiles})
		if err != nil {
			log.Fatalf("Unable to post new version of gist: %v", err)
		}

		log.Printf("Gist updated!")
	} else {
		log.Printf("No changes!")
	}
}

func main() {
	configuration := readConfig()
	myPublicIP := getMyIp(configuration)
	syncGistContent(configuration, myPublicIP)
}
