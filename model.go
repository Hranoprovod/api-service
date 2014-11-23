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

func getNode(c appengine.Context, slug string) *Node {
	var n Node
	if err := datastore.Get(c, getKey(c, slug), &n); err != nil {
		return nil
	}
	return &n
}

func saveNode(c appengine.Context, n Node) error {
	_, err := datastore.Put(c, getKey(c, n.Slug), &n)
	return err
}
