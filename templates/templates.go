package templates

import (
    "html/template"
    "log"
    "net/http"
)

var templates = map[string]*template.Template{}

func Init() {
  t, err := template.ParseFiles("templates/base.html")
  if err != nil {
      log.Fatal(err)
  }

  templates["main"] = t
}

func Render(name string, w http.ResponseWriter, c interface{}) {
    t, ok := templates[name]
    if !ok {
        log.Fatal("Invalid template: " + name)
    }

    t.Execute(w, c)
}
