package routes

import (
	"github.com/dimfeld/httptreemux"
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/urfave/negroni"
	"net/http"
	"os"
	"shortener/routes/common"
	"shortener/routes/dashboard"
	"shortener/routes/login"
	"shortener/routes/redirect"
)

func NewRouter() *negroni.Negroni {
	n := negroni.Classic()

	store := cookiestore.New([]byte(os.Getenv("S_SESSION_KEY")))
	n.Use(sessions.Sessions(os.Getenv("S_SESSION_NAME"), store))
	n.UseFunc(common.WithUser)

	r := httptreemux.NewContextMux()

	// Authenticated
	an := negroni.New(negroni.HandlerFunc(common.RequiresUser))

	r.Handler(http.MethodGet, common.PathDashboard, an.With(negroni.WrapFunc(dashboard.Dashboard)))

	r.Handler(http.MethodPost, common.PathLink, an.With(negroni.WrapFunc(dashboard.ShortenLink)))
	r.Handler(http.MethodDelete, common.PathLink, an.With(negroni.WrapFunc(dashboard.DeleteLink)))

	r.Handler(http.MethodGet, common.PathLogout, an.With(negroni.WrapFunc(login.Logout)))

	// Public
	fn := negroni.New(negroni.HandlerFunc(common.AutoDashboard))

	loginHandler := fn.With(negroni.WrapFunc(login.Login))
	r.Handler(http.MethodGet, common.PathLogin, loginHandler)
	r.Handler(http.MethodPost, common.PathLogin, loginHandler)

	registerHandler := fn.With(negroni.WrapFunc(login.Register))
	r.Handler(http.MethodGet, common.PathRegister, registerHandler)
	r.Handler(http.MethodPost, common.PathRegister, registerHandler)

	r.GET("/favicon.ico", http.NotFound)
	r.GET("/robots.txt", http.NotFound)

	r.GET("/:id", redirect.Redirect)

	n.UseHandler(r)

	return n
}
