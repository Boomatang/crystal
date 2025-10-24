package workflow

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/boomatang/crystal/pkg/logger"
	"github.com/jinzhu/copier"
)

func NewNode(event AdmissionRequest) *Node {
	return &Node{Kind: event.Kind.Kind, Name: event.Name, Data: event.Object}
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

type Direction int

const (
	Both Direction = iota
	Up
	Down
)

type NodeOptList struct {
	Direction Direction
}

type Node struct {
	Kind     string
	Name     string
	Data     any
	Parents  []*Node
	Childern []*Node
}

func (n *Node) HasChild(child *Node) bool {
	return slices.Contains(n.Childern, child)
}

func (n *Node) AddChild(child *Node) {
	if n.Childern == nil {
		n.Childern = []*Node{child}
		return
	}
	n.Childern = append(n.Childern, child)
}

func (n *Node) HasParent(parent *Node) bool {
	return slices.Contains(n.Parents, parent)
}

func (n *Node) AddParent(parent *Node) {
	if n.Parents == nil {
		n.Parents = []*Node{parent}
		return
	}
	n.Parents = append(n.Parents, parent)
}

func (n Node) Object(obj any) error {
	if dataMap, ok := n.Data.(map[string]any); ok {
		// Convert map to JSON and then to Deployment
		jsonData, err := json.Marshal(dataMap)
		if err != nil {
			return err

		}
		if err := json.Unmarshal(jsonData, obj); err != nil {
			return err
		}
	}
	return nil
}

func (n Node) String() string {
	return fmt.Sprintf("%s_%s", n.Kind, n.Name)
}

type NodeList struct {
	nodes   []*Node
	linkers []Link
}

func (nl *NodeList) Len() int {
	if nl.nodes == nil {
		return 0
	}
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

func (nl *NodeList) GetSubGraph(node *Node, opt NodeOptList) (*NodeList, error) {
	nodeList := &NodeList{}
	nodeCopy := &Node{}

	copier.Copy(nodeCopy, node)
	if opt.Direction != Up && opt.Direction != Both {
		nodeCopy.Parents = make([]*Node, 0)
	}

	if opt.Direction != Down && opt.Direction != Both {
		nodeCopy.Childern = make([]*Node, 0)
	}

	nodeList.Add(nodeCopy)
	if opt.Direction == Up || opt.Direction == Both {
		followParent(nodeList, nodeCopy)
	}

	if opt.Direction == Down || opt.Direction == Both {
		followChild(nodeList, nodeCopy)
	}

	return nodeList, nil
}

func followParent(nl *NodeList, node *Node) {
	logger.Log.Debug("traversing parent nodes", "node", node.String())
	for _, n := range node.Parents {
		nodeCopy := &Node{}
		copier.Copy(nodeCopy, n)
		nodeCopy.Childern = make([]*Node, 0)
		nl.Add(nodeCopy)
		followParent(nl, nodeCopy)
	}
}

func followChild(nl *NodeList, node *Node) {
	logger.Log.Debug("traversing child nodes", "node", node.String())
	for _, n := range node.Childern {
		nodeCopy := &Node{}
		copier.Copy(nodeCopy, n)
		nodeCopy.Parents = make([]*Node, 0)
		nl.Add(nodeCopy)
		followChild(nl, nodeCopy)
	}
}

func (nl *NodeList) Add(node *Node) {
	nl.nodes = append(nl.nodes, node)
}

func (nl *NodeList) GetNodes(t string) []*Node {
	// WARNING:this may need to be returning a NodeList.
	// If so then I think the NodeList needs a refactor also.
	nodes := make([]*Node, 0)
	for _, n := range nl.nodes {
		if n.Kind == t {
			nodes = append(nodes, n)
		}
	}
	return nodes
}

func (nl *NodeList) Link(node *Node) {
	// FIXME: this should be able to be done in one pass over the list of nodes
	kind := node.Kind
	pLinkers := make([]Link, 0)
	cLinkers := make([]Link, 0)
	for _, l := range nl.linkers {
		if l.Child == kind {
			pLinkers = append(pLinkers, l)
		}
		if l.Parent == kind {
			cLinkers = append(cLinkers, l)
		}
	}

	for _, linker := range pLinkers {
		for _, n := range nl.nodes {
			if n.Kind == linker.Parent && !n.HasChild(node) {
				if linker.LinkFunc(n, node) {
					n.AddChild(node)
					node.AddParent(n)
				}
			}
		}
	}

	for _, linker := range cLinkers {
		for _, n := range nl.nodes {
			if n.Kind == linker.Child && !n.HasParent(node) {
				if linker.LinkFunc(n, node) {
					n.AddParent(node)
					node.AddChild(n)
				}
			}
		}
	}
}

func (nl *NodeList) SetLinker(link Link) {
	nl.linkers = append(nl.linkers, link)
}

func (nl *NodeList) Render() string {
	s := ""

	seen := make([]*Node, nl.Len())

	for _, n := range nl.nodes {
		for _, p := range n.Parents {
			pString := fmt.Sprintf("%s -> %s", p, n)
			if !slices.Contains(seen, p) {
				s = fmt.Sprintf("%v\n%v", s, pString)
			}
		}
		for _, c := range n.Childern {
			cString := fmt.Sprintf("%s -> %s", n, c)
			if !slices.Contains(seen, c) {
				s = fmt.Sprintf("%v\n%v", s, cString)
			}
		}
		s = fmt.Sprintf("%v\n%s", s, n)

		seen = append(seen, n)
	}

	s = fmt.Sprintf("digraph {%v\n}", s)

	logger.Log.Debug("graph render output", "length", len(s))

	return s
}
