package apiservice

import (
	"appengine"
	"appengine/datastore"
	"appengine/search"
	"time"
)

type Node struct {
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	Calories     float64   `json:"calories"`
	Fat          float64   `json:"fat"`
	Carbohydrate float64   `json:"carbohydrate"`
	Protein      float64   `json:"protein"`
	Barcode      string    `json:"barcode,omitempty"`
	UserId       string    `json:"-"`
	Created      time.Time `json:"created"`
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

func saveNodeSearch(c appengine.Context, n *Node) error {
	index, err := search.Open("Node")
	if err != nil {
		return err
	}
	_, err = index.Put(c, n.Slug, n)
	if err != nil {
		return err
	}
	return nil
}

func searchNodes(c appengine.Context, q string, limit int) []Node {
	var nodes []Node
	index, _ := search.Open("Node")
	for t := index.Search(c, q, nil); ; {
		limit--
		if limit == 0 {
			break
		}
		var node Node
		_, err := t.Next(&node)
		if err == search.Done {
			break
		}
		if err != nil {
			break
		}
		nodes = append(nodes, node)
	}
	return nodes
}
