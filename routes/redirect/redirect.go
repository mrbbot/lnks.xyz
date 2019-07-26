package redirect

import (
	"github.com/dimfeld/httptreemux"
	"github.com/go-redis/redis"
	"log"
	"net/http"
	"shortener/db"
	"shortener/routes/common"
	"shortener/routes/templates"
	"time"
)

type errorData struct {
	Flashes []interface{}
	Error   string
}

func (d *errorData) SetFlashes(flashes []interface{}) {
	d.Flashes = flashes
}

const (
	errUnexpected = "An unexpected error occurred"
	errNotFound   = "Not found"
)

func Redirect(w http.ResponseWriter, r *http.Request) {
	params := httptreemux.ContextParams(r.Context())
	id := params["id"]

	linkKey := common.RedisLinkNamespace + r.Host + ":" + id

	url, err := db.R.HGet(linkKey, "url").Result()
	if err == redis.Nil {
		w.WriteHeader(http.StatusNotFound)
		templates.Render(w, r, "error", &errorData{
			Error: errNotFound,
		})
		return
	} else if err != nil {
		log.Printf("err getting url for id \"%s\": %v", id, err)
		w.WriteHeader(http.StatusInternalServerError)
		templates.Render(w, r, "error", &errorData{
			Error: errUnexpected,
		})
		return
	}

	go func() {
		_, err := db.R.HIncrBy(linkKey, "clicks", 1).Result()
		if err != nil {
			log.Printf("err incrementing click count for id \"%s\": %v", id, err)
		}

		_, err = db.R.HSet(linkKey, "lastClicked", time.Now().Format(common.LastClickLayout)).Result()
		if err != nil {
			log.Printf("err setting last click time for id \"%s\": %v", id, err)
		}
	}()

	http.Redirect(w, r, url, http.StatusFound)
}
