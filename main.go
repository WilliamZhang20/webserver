package main

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"os"
	"regexp"

	"webserver/data"
	"webserver/templates"
)

var validPath = regexp.MustCompile(("^/(edit|save|view)/([a-zA-Z0-9]+)$"))
var pageStore data.PageStore

type Page struct {
	Title string // a string title
	Body  []byte // slice of bytes
}

func (p *Page) save() error { // writes bytes to a text file
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("invalid Page Title")
	}
	return m[2], nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *data.Page) {
	var buf *bytes.Buffer
	var err error
	if buf, err = templates.ExecuteTemplate(tmpl, p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Printf("Error writing template to response: %v", err)
	}
}

// To avoid repetitive functionality, let's make a wrapper function
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2]) // call closure
	} // return closure
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := pageStore.LoadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := pageStore.LoadPage(title)
	if err != nil {
		p = &data.Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &data.Page{Title: title, Body: []byte(body)}
	err := pageStore.SavePage(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
	templates.LoadTemplates()
	dataSourceName := "user=postgres password=william_postgresql dbname=webserver_data sslmode=disable"
	store, err := data.NewPostgresPageStore(dataSourceName)
	if err != nil {
		log.Fatalf("Failed to initialize page store: %v", err)
	}
	pageStore = store
	// One can see the use of the Strategy Design Pattern with makeHandler calls...
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
