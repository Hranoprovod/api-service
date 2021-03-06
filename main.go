package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	. "github.com/gorilla/feeds"
	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
	"github.com/sourcegraph/sitemap"

	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
)

const (
	pageTitle = "Hranoprovod"
)

type TemplateData map[string]interface{}

var helperFuncs = template.FuncMap{
	"valToStr":  valToStr,
	"timeToStr": timeToStr,
}

func NewData(title string, description string) TemplateData {
	data := make(TemplateData)
	data["Title"] = title
	data["Description"] = description
	data["Heading"] = title
	return data
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/item/{slug}", itemHandler)
	r.HandleFunc("/add", addHandler)
	r.HandleFunc("/save", saveHandler)
	r.HandleFunc("/feed", feedHandler)
	r.HandleFunc("/search", searchHandler)
	r.HandleFunc("/sitemap.xml", sitemapHandler)
	r.HandleFunc("/api/v1/search", apiSearchHandler)
	http.Handle("/", r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func render(data TemplateData, w http.ResponseWriter, r *http.Request, filenames ...string) {
	c := appengine.NewContext(r)
	data["User"] = user.Current(c)
	url, _ := user.LoginURL(c, "/")
	data["LoginURL"] = url
	url, _ = user.LogoutURL(c, "/")
	data["LogoutURL"] = url
	t := template.New("layout.html")
	t.Funcs(helperFuncs)
	if err := template.Must(t.ParseFiles(filenames...)).Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	data := NewData(pageTitle, pageTitle)
	data["Latest"] = getLatestNodes(r.Context(), 10)
	render(data, w, r, "templates/layout.html", "templates/index.html")
}

func itemHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	vars := mux.Vars(r)
	node := getNode(c, vars["slug"])
	if node == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	data := NewData((*node).Name+" - "+pageTitle, (*node).Name)
	data["Node"] = *node
	data["Heading"] = (*node).Name
	render(data, w, r, "templates/layout.html", "templates/item.html")
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	render(NewData("Add food", "Add new tood to database"), w, r, "templates/layout.html", "templates/add.html")
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}
	n := Node{
		Name:         r.FormValue("name"),
		Slug:         slug.Make(r.FormValue("name")),
		Calories:     getFloat(r.FormValue("calories")),
		Fat:          getFloat(r.FormValue("fat")),
		Carbohydrate: getFloat(r.FormValue("carbohydrate")),
		Protein:      getFloat(r.FormValue("protein")),
		Barcode:      r.FormValue("barcode"),
		UserId:       u.ID,
		Created:      time.Now(),
	}
	err := saveNode(c, n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = saveNodeSearch(c, &n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/item/"+n.Slug, http.StatusFound)
}

func feedHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	nodes := getLatestNodes(c, 20)

	feed := &Feed{
		Title:       pageTitle,
		Link:        &Link{Href: "http://" + r.Host + "/"},
		Description: "Nutrition information",
		Author:      &Author{"Evgeniy Vasilev", "aquilax@gmail.com"},
		Created:     time.Now(),
	}
	for _, node := range nodes {
		feed.Items = append(feed.Items, &Item{
			Title:   node.Name,
			Link:    &Link{Href: "http://" + r.Host + "/item/" + node.Slug},
			Created: node.Created,
		})
	}
	rss, err := feed.ToRss()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/rss+xml")
	fmt.Fprint(w, rss)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if q == "" {
		http.Error(w, "No query string found", http.StatusBadRequest)
		return
	}
	data := NewData(pageTitle, pageTitle)
	data["Heading"] = q
	data["SearchString"] = q
	data["Results"] = searchNodes(c, q, 0)
	render(data, w, r, "templates/layout.html", "templates/search.html")
}

func sitemapHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	latest := getLatestNodes(c, sitemap.MaxURLs)
	var urlSet sitemap.URLSet
	for _, node := range latest {
		urlSet.URLs = append(urlSet.URLs, sitemap.URL{
			Loc:        "http://" + r.Host + "/item/" + node.Slug,
			LastMod:    &node.Created,
			ChangeFreq: sitemap.Never,
			Priority:   0.7,
		})
	}
	xml, err := sitemap.Marshal(&urlSet)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/xml")
	w.Write(xml)
}
