package main

import (
	"context"

	"github.com/Hranoprovod/shared"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/search"
)

type Node shared.APINode

func getKey(c context.Context, slug string) *datastore.Key {
	return datastore.NewKey(c, "Node", slug, 0, nil)
}

func getLatestNodes(c context.Context, limit int) []Node {
	q := datastore.NewQuery("Node").Order("-Created").Limit(limit)
	var nodes []Node
	q.GetAll(c, &nodes)
	return nodes
}

func getNode(c context.Context, slug string) *Node {
	var n Node
	if err := datastore.Get(c, getKey(c, slug), &n); err != nil {
		return nil
	}
	return &n
}

func saveNode(c context.Context, n Node) error {
	_, err := datastore.Put(c, getKey(c, n.Slug), &n)
	return err
}

func saveNodeSearch(c context.Context, n *Node) error {
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

func searchNodes(c context.Context, q string, limit int) []Node {
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
