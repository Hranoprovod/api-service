package apiservice

import (
	"appengine"
	"appengine/datastore"
	"appengine/search"
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

type SearchNode struct {
	Name string
	Slug string
	Calories     float64
	Fat          float64
	Carbohydrate float64
	Protein      float64
	Barcode string
	Created      time.Time
}

func (n Node) NewSearchNode() *SearchNode {
	return &SearchNode{
		Name : n.Name,
		Slug: n.Slug,
		Calories: intToFloat64(n.Calories),
		Fat: intToFloat64(n.Fat),
		Carbohydrate: intToFloat64(n.Carbohydrate),
		Protein: intToFloat64(n.Protein),
		Barcode: n.Barcode,
		Created: n.Created,
	}
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

func saveNodeSearch(c appengine.Context, sn *SearchNode) error {
	index, err := search.Open("SearchNode")
	if err != nil {
		return err
	}
	_, err = index.Put(c, sn.Slug, sn)
	if err != nil {
		return err
	}
	return nil
}

func searchNodes(c appengine.Context, q string, limit int) []SearchNode {
	var nodes []SearchNode 
	index, _ := search.Open("SearchNode")
	for t := index.Search(c, q, nil); ; {
		limit--
		if (limit == 0) {
        	break
        }
        var node SearchNode
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