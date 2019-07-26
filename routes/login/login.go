package login

import (
	"github.com/go-redis/redis"
	"github.com/goincremental/negroni-sessions"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"shortener/db"
	"shortener/routes/common"
	"shortener/routes/templates"
)

type loginTemplateData struct {
	Flashes     []interface{}
	Username    string
	Registering bool
}

func (d *loginTemplateData) SetFlashes(flashes []interface{}) {
	d.Flashes = flashes
}

const (
	bcryptCost = 14

	msgUnexpectedError       = "An unexpected error occurred!"
	msgInvalidCredentials    = "Invalid username or password!"
	msgUsernamePasswordEmpty = "You must enter a username and password!"
	msgLoggedIn              = "You are now logged in"
)

func Login(w http.ResponseWriter, r *http.Request) {
	s := sessions.GetSession(r)
	var username string

	if r.Method == http.MethodPost {
		username = r.PostFormValue("username")
		password := r.PostFormValue("password")

		if username == "" || password == "" {
			s.AddFlash(msgUsernamePasswordEmpty)
		} else {
			hashedPassword, err := db.R.Get(common.RedisPasswordNamespace + username).Result()
			if err == redis.Nil {
				log.Print("err in login: could not find user")
				s.AddFlash(msgInvalidCredentials)
			} else if err != nil {
				log.Printf("err getting password: %v", err)
				s.AddFlash(msgUnexpectedError)
			} else if err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
				log.Print("err in login: passwords do not match")
				s.AddFlash(msgInvalidCredentials)
			} else {
				validHost, err := db.R.SIsMember(common.RedisHostsNamespace+username, r.Host).Result()
				if err != nil {
					log.Printf("err checking valid host: %v", err)
					s.AddFlash(msgUnexpectedError)
				} else if !validHost {
					s.AddFlash(common.MsgInvalidHost)
				} else {
					s.Set(common.UserKey, username)
					s.AddFlash(msgLoggedIn)
					http.Redirect(w, r, common.PathDashboard, http.StatusFound)
					return
				}
			}
		}
	}

	templates.Render(w, r, "login", &loginTemplateData{
		Username: username,
	})
}
