package templates

import (
    "html/template"
    "log"
    "net/http"
)

var templates = map[string]*template.Template{}

func initTemplate(name string, filenames ...string) {
    t, err := template.ParseFiles(filenames...)
    if err != nil {
        log.Fatal("Error retrieving or compiling the " + name + " template.")
    }

    templates[name] = t
}

func Init() {
  initTemplate("leaderboard", "templates/base.html")
}

func Render(name string, w http.ResponseWriter, c interface{}) {
    t, ok := templates[name]
    if !ok {
        log.Fatal("Invalid template: " + name)
    }

    t.Execute(w, c)
}
