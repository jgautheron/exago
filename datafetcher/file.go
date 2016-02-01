package datafetcher

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/jgautheron/exago-service/config"
)

func GetFileContents(log *log.Entry, repository, path string) ([]byte, error) {
	// Only github implemented for now
	return fromGithub(log, repository, path)
}

func fromGithub(log *log.Entry, repository, path string) ([]byte, error) {
	sp := strings.Split(repository, "/")
	username, repo := sp[1], sp[2]

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s?access_token=%s", username, repo, path, config.Get("GithubAccessToken"))
	log.Info(url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err = json.Unmarshal(raw, &data); err != nil {
		return nil, err
	}

	return base64.StdEncoding.DecodeString(data["content"].(string))
}
