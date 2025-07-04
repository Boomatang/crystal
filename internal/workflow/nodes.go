package workflow

import (
	"fmt"
	"slices"
)

func NewNode(event Event) *Node {
	return &Node{Kind: event.Kind, Name: event.Name}
}

func NewNodeList() *NodeList {
	nodes := make([]*Node, 0)
	linkers := make([]Link, 0)
	return &NodeList{nodes: nodes, linkers: linkers}
}

type Link struct {
	Parent   string
	Child    string
	LinkFunc func(*Node, *Node) bool
}

type Node struct {
	Kind     string
	Name     string
	Parents  []*Node
	Childern []*Node
}

func (n *Node) AddChild(child *Node) {
	if n.Childern == nil {
		n.Childern = []*Node{child}
		return
	}
	n.Childern = append(n.Childern, child)
}

func (n *Node) AddParent(parent *Node) {
	if n.Parents == nil {
		n.Parents = []*Node{parent}
		return
	}
	n.Parents = append(n.Parents, parent)
}

func (n Node) String() string {
	return fmt.Sprintf("%s_%s", n.Kind, n.Name)
}

type NodeList struct {
	nodes   []*Node
	linkers []Link
}

func (nl *NodeList) Len() int {
	return len(nl.nodes)
}

func (nl *NodeList) Contains(node Node) bool {
	for _, n := range nl.nodes {
		if n.Kind == node.Kind && n.Name == node.Name {
			return true
		}
	}
	return false
}

func (nl *NodeList) Get(kind string, name string) *Node {
	for _, n := range nl.nodes {
		if n.Kind == kind && n.Name == name {
			return n
		}
	}
	return nil

}

func (nl *NodeList) Add(node *Node) {
	nl.nodes = append(nl.nodes, node)
}

func (nl *NodeList) Link(node *Node) {
	// FIXME: this should be able to be done in one pass over the list of nodes
	kind := node.Kind
	pLinkers := make([]Link, 0)
	for _, l := range nl.linkers {
		if l.Child == kind {
			pLinkers = append(pLinkers, l)
		}
	}

	for _, linker := range pLinkers {
		for _, n := range nl.nodes {
			if n.Kind == linker.Parent && linker.LinkFunc(n, node) {
				n.AddChild(node)
				node.AddParent(n)
			}
		}
	}

	cLinkers := make([]Link, 0)
	for _, l := range nl.linkers {
		if l.Parent == kind {
			cLinkers = append(cLinkers, l)
		}
	}

	for _, linker := range cLinkers {
		for _, n := range nl.nodes {
			if n.Kind == linker.Child && linker.LinkFunc(n, node) {
				n.AddParent(node)
				node.AddChild(n)
			}
		}
	}
}

func (nl *NodeList) SetLinker(link Link) {
	nl.linkers = append(nl.linkers, link)
}

func (nl *NodeList) Render() string {
	s := ""

	seen := make([]*Node, len(nl.nodes))

	for _, n := range nl.nodes {
		for _, p := range n.Parents {
			if !slices.Contains(seen, p) {
				s = fmt.Sprintf("%v\n%s -> %s", s, p, n)
			}
		}
		for _, c := range n.Childern {
			if !slices.Contains(seen, c) {
				s = fmt.Sprintf("%v\n%s -> %s", s, n, c)
			}
		}

		seen = append(seen, n)
	}

	s = fmt.Sprintf("digraph {%v\n}", s)

	return s
}
