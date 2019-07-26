package common

import (
	"context"
	"fmt"
	"github.com/goincremental/negroni-sessions"
	"net/http"
	"shortener/db"
)

const (
	UserKey = "user"

	msgRequiresLogin = "You must be logged in to access this page!"
	MsgInvalidHost   = "Your account does not have access to this domain!"
)

func WithUser(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	s := sessions.GetSession(r)
	user := s.Get(UserKey)
	ctx := context.WithValue(r.Context(), UserKey, user)
	r = r.WithContext(ctx)
	next(w, r)
}

func GetUser(r *http.Request) string {
	if user, ok := r.Context().Value(UserKey).(string); ok {
		return user
	}
	return ""
}

func RequiresUser(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	s := sessions.GetSession(r)
	var user string

	if user = GetUser(r); user == "" {
		s.AddFlash(msgRequiresLogin)
		http.Redirect(w, r, PathLogin, http.StatusFound)
	} else {
		validHost, err := db.R.SIsMember(RedisHostsNamespace+user, r.Host).Result()
		if err != nil {
			http.Error(w, fmt.Sprintf("err checking valid host: %v", err), http.StatusInternalServerError)
		} else if !validHost {
			s.AddFlash(MsgInvalidHost)
			http.Redirect(w, r, PathLogin, http.StatusFound)
		} else {
			next(w, r)
		}
	}
}

func AutoDashboard(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if user := GetUser(r); user == "" {
		next(w, r)
	} else {
		validHost, err := db.R.SIsMember(RedisHostsNamespace+user, r.Host).Result()
		if err != nil {
			http.Error(w, fmt.Sprintf("err checking valid host: %v", err), http.StatusInternalServerError)
		} else if !validHost {
			next(w, r)
		} else {
			http.Redirect(w, r, PathDashboard, http.StatusFound)
		}
	}
}
