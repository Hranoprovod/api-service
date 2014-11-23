package apiservice

import (
	"appengine"
	"appengine/user"
	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
	"html/template"
	"net/http"
	"time"
)

type TemplateData map[string]interface{}

var helperFuncs = template.FuncMap{
	"valToStr":  valToStr,
	"timeToStr": timeToStr,
}

func NewData() TemplateData {
	return make(TemplateData)
}

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/item/{itemID}", itemHandler)
	r.HandleFunc("/add", addHandler)
	r.HandleFunc("/save", saveHandler)
	http.Handle("/", r)
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
	c := appengine.NewContext(r)
	data := make(map[string]interface{})
	data["Latest"] = getLatestNodes(c, 10)
	render(data, w, r, "templates/layout.html", "templates/index.html")
}

func itemHandler(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	render(data, w, r, "templates/layout.html", "templates/item.html")
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	render(NewData(), w, r, "templates/layout.html", "templates/add.html")
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
		Calories:     floatToInt(getFloat(r.FormValue("calories"))),
		Fat:          floatToInt(getFloat(r.FormValue("fat"))),
		Carbohydrate: floatToInt(getFloat(r.FormValue("carbohydrate"))),
		Protein:      floatToInt(getFloat(r.FormValue("protein"))),
		Barcode:      r.FormValue("barcode"),
		UserId:       u.ID,
		Created:      time.Now(),
	}
	err := saveNode(c, n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
