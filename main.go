package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/alecthomas/kingpin"
)

type gitlabPushEvent struct {
	Repository struct {
		Slug       string `json:"slug"`         // bitbucket
		GitHTTPURL string `json:"git_http_url"` // gitlab
	} `json:"repository"`
}

var (
	addr     string
	server   string
	user     string
	password string
)

func main() {
	kingpin.Flag("addr", "Listen address").Default(":8080").StringVar(&addr)
	kingpin.Flag("server", "TeamCity server URL").Default("http://localhost:8111").StringVar(&server)
	kingpin.Flag("user", "TeamCity user").Required().StringVar(&user)
	kingpin.Flag("password", "TeamCity password").Required().StringVar(&password)
	kingpin.Parse()

	http.HandleFunc("/", handlePushEvent)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("HTTP listen:", err)
	}
}

func handlePushEvent(w http.ResponseWriter, req *http.Request) {
	var ev gitlabPushEvent
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&ev); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	tcHook(ev)
}

func tcHook(ev gitlabPushEvent) {
	repo := ev.Repository.Slug
	if ev.Repository.GitHTTPURL != "" {
		repourl, err := url.Parse(ev.Repository.GitHTTPURL)
		if err != nil {
			log.Println("Un-parseable Git HTTP URL:", ev.Repository.GitHTTPURL)
			return
		}

		repo = repourl.Path
		repo = strings.TrimSuffix(repo, path.Ext(repo))
	}

	locator := fmt.Sprintf("vcsRoot:(type:jetbrains.git,property:(name:url,value:%s,matchType:contains,ignoreCase:true),count:99999),count:99999", repo)
	values := url.Values{"locator": []string{locator}}

	reqpath := server + "/app/rest/vcs-root-instances/commitHookNotification?" + values.Encode()
	req, err := http.NewRequest(http.MethodPost, reqpath, nil)
	if err != nil {
		log.Println("Invalid resulting request:", reqpath)
		return
	}
	req.SetBasicAuth(user, password)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("HTTP error:", err)
		return
	}

	bs, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Println("HTTP error:", err)
		return
	}

	if resp.StatusCode != http.StatusAccepted {
		log.Println("HTTP error:", resp.Status)
		return
	}

	bs = bytes.TrimSpace(bs)
	log.Printf("Delivered hook for %s: %s (%s)", repo, bs, resp.Status)
}
