package login

import (
	"github.com/goincremental/negroni-sessions"
	"net/http"
	"shortener/routes/common"
)

const (
	msgLoggedOut = "You are now logged out"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	s := sessions.GetSession(r)
	s.Set(common.UserKey, "")
	s.AddFlash(msgLoggedOut)
	http.Redirect(w, r, common.PathLogin, http.StatusFound)
}
