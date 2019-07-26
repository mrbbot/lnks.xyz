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

const (
	msgCodeUsernamePasswordEmpty = "You must enter a registration code, username and password!"
	msgInvalidRegistrationCode   = "Invalid registration code!"
	msgUsernameExists            = "A user already exists with that username!"
)

func Register(w http.ResponseWriter, r *http.Request) {
	s := sessions.GetSession(r)
	var username string

	if r.Method == http.MethodPost {
		code := r.PostFormValue("code")
		username = r.PostFormValue("username")
		password := r.PostFormValue("password")

		if code == "" || username == "" || password == "" {
			s.AddFlash(msgCodeUsernamePasswordEmpty)
		} else {
			allowedHost, err := db.R.Get(common.RedisRegistrationCodeNamespace + code).Result()

			if err == redis.Nil {
				log.Print("err in register: cannot find code")
				s.AddFlash(msgInvalidRegistrationCode)
			} else if err != nil {
				log.Printf("err getting code: %v", err)
				s.AddFlash(msgUnexpectedError)
			} else if allowedHost != r.Host {
				log.Printf(
					"err in register: code '%s' not valid for host '%s', current host '%s'",
					code, allowedHost, r.Host,
				)
				s.AddFlash(msgInvalidRegistrationCode)
			} else {
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)

				if err != nil {
					log.Printf("err hashing password: %v", err)
					s.AddFlash(msgUnexpectedError)
				} else {
					set, err := db.R.SetNX(common.RedisPasswordNamespace+username, hashedPassword, 0).Result()

					if err != nil {
						log.Printf("err storing password: %v", err)
						s.AddFlash(msgUnexpectedError)
					} else if !set {
						s.AddFlash(msgUsernameExists)
					} else {
						_, err := db.R.SAdd(common.RedisHostsNamespace+username, allowedHost).Result()

						if err != nil {
							log.Printf("err storing allowed host: %v", err)
							s.AddFlash(msgUnexpectedError)
						} else {
							s.Set(common.UserKey, username)
							s.AddFlash(msgLoggedIn)
							http.Redirect(w, r, common.PathDashboard, http.StatusFound)
							return
						}
					}
				}
			}
		}
	}

	templates.Render(w, r, "login", &loginTemplateData{
		Username:    username,
		Registering: true,
	})
}
