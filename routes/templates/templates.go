package templates

import (
	"fmt"
	sessions "github.com/goincremental/negroni-sessions"
	"html/template"
	"net/http"
	"os"
)

var tmpl *template.Template

type FlashStorer interface {
	SetFlashes(flashes []interface{})
}

func Render(w http.ResponseWriter, r *http.Request, name string, data FlashStorer) {
	s := sessions.GetSession(r)
	data.SetFlashes(s.Flashes())

	var err error
	if tmpl == nil || os.Getenv("GO_ENV") == "development" {
		tmpl, err = template.ParseGlob("templates/**/*.html")
		if err != nil {
			http.Error(w, fmt.Sprintf("err parsing template %s: %v", name, err), http.StatusInternalServerError)
			return
		}
	}

	err = tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("err executing template %s: %v", name, err), http.StatusInternalServerError)
	}
}
