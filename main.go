package apiservice

import (
	"time"
	"strconv"
	"net/http"
	"appengine"
	"appengine/user"
	"appengine/datastore"
	"html/template"
)

const(
	PRECISION = 1000
)

func init() {
	http.HandleFunc("/", index)
	http.HandleFunc("/add", add)
	http.HandleFunc("/save", save)
}

func getLatestNodes(c appengine.Context, limit int) []Node {
	q := datastore.NewQuery("Node").Order("-Created").Limit(limit)
	var nodes []Node
	q.GetAll(c, &nodes)
	return nodes
}

func index(w http.ResponseWriter, r *http.Request) {
	var indexTmpl = template.Must(template.ParseFiles("templates/layout.html", "templates/index.html"))
	c := appengine.NewContext(r)
	data := make(map[string]interface{})
	data["Latest"] = getLatestNodes(c, 10);
	if err := indexTmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func add(w http.ResponseWriter, r *http.Request) {
	var addTmpl = template.Must(template.ParseFiles("templates/layout.html", "templates/add.html"))
	if err := addTmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getFloat(value string) float32 {
	num, _ := strconv.ParseFloat(value, 32)
	return float32(num)
}

func floatToInt(value float32) int {
	return int(value * PRECISION) 
}

type Node struct {
	Name string
	Calories int
	Fat int
	Carbohydrate int
	Protein int
	Barcode string
	UserId string
	Created time.Time
}

func save(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	userId := ""
	if u != nil {
		userId = u.ID
	}
	n := Node{
		Name: r.FormValue("name"),
		Calories: floatToInt(getFloat(r.FormValue("calories"))),
		Fat: floatToInt(getFloat(r.FormValue("fat"))),
		Carbohydrate: floatToInt(getFloat(r.FormValue("carbohydrate"))),
		Protein: floatToInt(getFloat(r.FormValue("protein"))),
		Barcode: r.FormValue("barcode"),
		UserId: userId,
		Created:    time.Now(),
	}
	key := datastore.NewKey(c, "Node", n.Name, 0, nil)
	_, err := datastore.Put(c, key, &n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}


