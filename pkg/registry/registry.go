package registry

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/serhio83/druid/pkg/config"
	"github.com/serhio83/druid/pkg/utils"
)

const (
	drCatalogURL = "/v2/_catalog?n=1000"
	drTagsList   = "/tags/list?n=1000"
	ver1         = "application/vnd.docker.distribution.manifest.v1+json"
	ver2         = "application/vnd.docker.distribution.manifest.v2+json"
)

// Tag name
type Tag string

// Repo name or registry image path
type Repo string

// Tags contains image tags list
type Tags struct {
	Name string
	Tags []Tag `json:"tags"`
}

// Repos ...
type Repos struct {
	Repositories []Repo `json:"repositories"`
}

// ImageDigest ...
type ImageDigest struct {
	Digest map[string]string
}

// Images ...
type Images struct {
	Name    string
	Digests ImageDigest
}

// V1Manifest ...
type V1Manifest struct {
	Name    string    `json:"name"`
	Tag     string    `json:"tag"`
	History []History `json:"history"`
}

// History ...
type History struct {
	V1Compatibility string `json:"v1Compatibility"`
}

func registryRequest(url, m string, c *config.Config, version string) *http.Response {
	config := &tls.Config{
		InsecureSkipVerify: true,
	}
	tr := &http.Transport{TLSClientConfig: config}
	client := &http.Client{Transport: tr}

	req, _ := http.NewRequest(m, url, nil)
	req.SetBasicAuth(c.RegUser, c.RegPass)
	req.Header.Set("Accept", version)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer client.CloseIdleConnections()
	return resp
}

// DelManifest ...
func DelManifest(c *config.Config, registry, manifest string) string {
	resp := registryRequest(c.RegHost+":"+c.RegPort+"/v2/"+registry+"/manifests/"+manifest, "DELETE", c, ver2)
	_, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	return resp.Status
}

func timeToUnix(ts string) string {
	t, _ := time.Parse(time.RFC3339, ts)
	return strconv.FormatInt(t.Unix(), 10)
}

func unixToTime(date string) *time.Time {
	i, err := strconv.ParseInt(date, 10, 64)
	if err != nil {
		log.Println(err)
	}
	tm := time.Unix(i, 0)
	return &tm
}

// GetCreationDate ...
func GetCreationDate(c *config.Config, registry, tagName string) string {

	resp := registryRequest(c.RegHost+":"+c.RegPort+"/v2/"+registry+"/manifests/"+tagName, "GET", c, ver1)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	var v1manifest V1Manifest
	parseErr := json.Unmarshal(body, &v1manifest)
	if parseErr != nil {
		log.Fatal(parseErr)
	}

	if len(v1manifest.History) > 0 {
		in := []byte(v1manifest.History[0].V1Compatibility)
		var raw map[string]interface{}
		json.Unmarshal(in, &raw)
		ct := timeToUnix(raw["created"].(string))
		return ct
	}
	log.Println(utils.Envelope(fmt.Sprintf("%s [time.search] %s:%s No creation time found", logHeader, registry, tagName)))
	return "nodate"
}

// GetManifest ...
func GetManifest(c *config.Config, registry, tagName string) (string, string) {
	resp := registryRequest(c.RegHost+":"+c.RegPort+"/v2/"+registry+"/manifests/"+tagName, "GET", c, ver2)
	_, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	return resp.Header.Get("Docker-Content-Digest"), resp.Status
}

// ListTags ...
func ListTags(c *config.Config, registry string) *Tags {
	resp := registryRequest(c.RegHost+":"+c.RegPort+"/v2/"+registry+drTagsList, "GET", c, ver2)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	var tags Tags
	parseErr := json.Unmarshal(body, &tags)
	if parseErr != nil {
		log.Fatal(parseErr)
	}
	return &tags
}

// ListImages ...
func ListImages(c *config.Config) (*Repos, error) {
	resp := registryRequest(c.RegHost+":"+c.RegPort+drCatalogURL, "GET", c, ver2)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var repos Repos
	parseErr := json.Unmarshal(body, &repos)
	if parseErr != nil {
		log.Fatal(parseErr)
		return nil, parseErr
	}
	return &repos, nil
}
