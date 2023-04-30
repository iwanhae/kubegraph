package graph

import (
	"reflect"
	"sort"
)

type Graph struct {
	nodes    map[string]Node
	backlink map[string]map[string]interface{}
	event    chan NodeEvent
}

func NewGraph() (*Graph, chan NodeEvent) {
	ch := make(chan NodeEvent, 0)
	return &Graph{
		nodes:    make(map[string]Node),
		backlink: make(map[string]map[string]interface{}),
		event:    ch,
	}, ch
}

type Node struct {
	ID    string   `json:"id"`
	Edges []string `json:"edges,omitempty"`

	Content any `json:"content"`
}

type NodeEvent struct {
	Status EventStatus
	ID     string
}

type EventStatus string

const (
	Created EventStatus = "created"
	Updated EventStatus = "updated"
	Deleted EventStatus = "deleted"
)

func (g *Graph) newEvent(s EventStatus, id string) {
	g.event <- NodeEvent{Status: s, ID: id}
}

func (g *Graph) GetNode(id string) *Node {
	if v, ok := g.nodes[id]; ok {
		return &v
	}
	return nil
}

func (g *Graph) setBackLink(from string, to string, activate bool) {
	backlinks, ok := g.backlink[from]
	if !ok {
		backlinks = make(map[string]interface{})
	}
	if activate {
		backlinks[to] = true
	} else {
		delete(backlinks, to)
	}
	g.backlink[from] = backlinks
}

func (g *Graph) getBacklinks(from string) []string {
	backlinks, ok := g.backlink[from]
	result := []string{}
	if ok {
		for k := range backlinks {
			result = append(result, k)
		}
	}
	return result
}

func (g *Graph) UpdateNode(id string, linkTo []string, content any) {
	sort.Strings(linkTo)
	old, ok := g.nodes[id]

	// Delete existing one
	if content == nil {
		delete(g.nodes, id)
		g.newEvent(Deleted, id)
		for _, edge := range old.Edges {
			g.setBackLink(edge, id, false)
		}
		for _, ref := range g.getBacklinks(id) {
			g.newEvent(Updated, ref)
		}
		return
	}

	node := Node{
		ID:      id,
		Edges:   linkTo,
		Content: content,
	}

	// Create a New Node
	if !ok {
		g.nodes[id] = node
		g.newEvent(Created, id)
		for _, edge := range linkTo {
			g.setBackLink(edge, id, true)
		}
		for _, ref := range g.getBacklinks(id) {
			g.newEvent(Updated, ref)
		}
		return
	}

	sameEdge := reflect.DeepEqual(old.Edges, linkTo)
	sameContent := reflect.DeepEqual(old.Content, content)
	if !sameEdge {
		for _, edge := range old.Edges {
			g.setBackLink(edge, old.ID, false)
		}
		for _, edge := range linkTo {
			g.setBackLink(edge, id, true)
		}
	}
	if !sameEdge || !sameContent {
		g.nodes[id] = node
		g.newEvent(Updated, id)
	}
}
