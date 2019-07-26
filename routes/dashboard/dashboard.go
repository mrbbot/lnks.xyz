package dashboard

import (
	"fmt"
	"log"
	"net/http"
	"shortener/db"
	"shortener/routes/common"
	"shortener/routes/templates"
)

type dashboardTemplateData struct {
	Flashes    []interface{}
	ShortLinks []*shortLink
}

func (d *dashboardTemplateData) SetFlashes(flashes []interface{}) {
	d.Flashes = flashes
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	user := common.GetUser(r)

	userLinkIds, err := db.R.LRange(common.RedisUserLinksNamespace+user+":"+r.Host, 0, -1).Result()
	if err != nil {
		log.Printf("err getting user's links: %v", err)
		http.Error(w, fmt.Sprintf("err getting user's links: %v", err), http.StatusInternalServerError)
		return
	}

	links := make([]*shortLink, 0)
	for _, linkId := range userLinkIds {
		linkMap, err := db.R.HGetAll(common.RedisLinkNamespace + r.Host + ":" + linkId).Result()
		if err != nil {
			log.Printf("err getting link map for id \"%s\": %v", linkId, err)
			http.Error(w, fmt.Sprintf("err getting link map for id \"%s\": %v", linkId, err), http.StatusInternalServerError)
			return
		}
		link := newLinkFromMap(linkMap)
		links = append(links, link)
	}

	templates.Render(w, r, "dashboard", &dashboardTemplateData{
		ShortLinks: links,
	})
}
