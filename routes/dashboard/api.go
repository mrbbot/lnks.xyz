package dashboard

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"shortener/db"
	"shortener/routes/common"
	"strconv"
	"strings"
	"time"
)

type shortLink struct {
	Host        string `json:"host"`
	Id          string `json:"id"`
	Url         string `json:"url"`
	Clicks      int    `json:"clicks"`
	LastClicked string `json:"lastClicked"`
	Created     string `json:"created"`
}

func (l *shortLink) toMap() map[string]interface{} {
	return map[string]interface{}{
		"host":        l.Host,
		"id":          l.Id,
		"url":         l.Url,
		"clicks":      l.Clicks,
		"lastClicked": l.LastClicked,
		"created":     l.Created,
	}
}

func newLinkFromMap(m map[string]string) *shortLink {
	clicks, err := strconv.Atoi(m["clicks"])
	if err != nil {
		panic(err)
	}
	return &shortLink{
		Host:        m["host"],
		Id:          m["id"],
		Url:         m["url"],
		Clicks:      clicks,
		LastClicked: m["lastClicked"],
		Created:     m["created"],
	}
}

type shortenLinkRequest struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

type deleteLinkRequest struct {
	Id string `json:"id"`
}

func parseRequest(r *http.Request, req interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, req)
	if err != nil {
		return err
	}
	return nil
}

func ShortenLink(w http.ResponseWriter, r *http.Request) {
	user := common.GetUser(r)

	var req shortenLinkRequest
	err := parseRequest(r, &req)
	if err != nil {
		http.Error(w, fmt.Sprintf("err parsing body: %v", err), http.StatusInternalServerError)
		return
	}

	link := &shortLink{
		Host:        r.Host,
		Id:          req.Id,
		Url:         req.Url,
		Clicks:      0,
		LastClicked: "",
		Created:     time.Now().In(common.LocationLondon).Format(common.CreatedLayout),
	}

	lowerLinkId := strings.ToLower(link.Id)
	linkKey := common.RedisLinkNamespace + link.Host + ":" + lowerLinkId

	exists, err := db.R.Exists(linkKey).Result()
	if err != nil {
		log.Printf("err checking link existence: %v", err)
		http.Error(w, fmt.Sprintf("err checking link existence: %v", err), http.StatusInternalServerError)
		return
	}

	if exists == 1 {
		http.Error(w, "id already exists", http.StatusConflict)
		return
	}

	_, err = db.R.HMSet(linkKey, link.toMap()).Result()
	if err != nil {
		log.Printf("err storing link: %v", err)
		http.Error(w, fmt.Sprintf("err storing link: %v", err), http.StatusInternalServerError)
		return
	}

	_, err = db.R.LPush(common.RedisUserLinksNamespace+user+":"+link.Host, lowerLinkId).Result()
	if err != nil {
		log.Printf("err adding link to user's list: %v", err)
		http.Error(w, fmt.Sprintf("err adding link to user's list: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	err = encoder.Encode(link)
	if err != nil {
		http.Error(w, fmt.Sprintf("err encoding res: %v", err), http.StatusInternalServerError)
	}
}

func DeleteLink(w http.ResponseWriter, r *http.Request) {
	user := common.GetUser(r)

	var req deleteLinkRequest
	err := parseRequest(r, &req)
	if err != nil {
		http.Error(w, fmt.Sprintf("err parsing body: %v", err), http.StatusInternalServerError)
		return
	}

	lowerLinkId := strings.ToLower(req.Id)
	linkKey := common.RedisLinkNamespace + r.Host + ":" + lowerLinkId

	removed, err := db.R.Del(linkKey).Result()
	if err != nil {
		log.Printf("err deleting link hash: %v", err)
		http.Error(w, fmt.Sprintf("err deleting link hash: %v", err), http.StatusInternalServerError)
		return
	} else if removed == 0 {
		http.Error(w, "couldn't find id", http.StatusNotFound)
		return
	}

	_, err = db.R.LRem(common.RedisUserLinksNamespace+user+":"+r.Host, 0, lowerLinkId).Result()
	if err != nil {
		log.Printf("err deleting link from user's links: %v", err)
		http.Error(w, fmt.Sprintf("err deleting link from user's links: %v", err), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write([]byte("ok"))
}
