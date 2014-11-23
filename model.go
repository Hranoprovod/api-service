package apiservice

import (
	"appengine"
	"appengine/datastore"
	"time"
)

type Node struct {
	Name         string
	Slug         string
	Calories     int
	Fat          int
	Carbohydrate int
	Protein      int
	Barcode      string
	UserId       string
	Created      time.Time
}

func getKey(c appengine.Context, slug string) *datastore.Key {
	return datastore.NewKey(c, "Node", slug, 0, nil)
}

func getLatestNodes(c appengine.Context, limit int) []Node {
	q := datastore.NewQuery("Node").Order("-Created").Limit(limit)
	var nodes []Node
	q.GetAll(c, &nodes)
	return nodes
}

func getNode(c appengine.Context, slug string) Node {
	q := datastore.NewQuery("Node").Filter("Slug", slug)
	var nodes []Node
	q.GetAll(c, &nodes)
	return nodes[0]
}

func saveNode(c appengine.Context, n Node) error {
	_, err := datastore.Put(c, getKey(c, n.Slug), &n)
	return err
}
