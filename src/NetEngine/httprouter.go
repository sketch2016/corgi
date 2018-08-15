package NetEngine

import (
	"strings"
)

const (
	NODE_TYPE_COMMAND = iota
	NODE_TYPE_PARAM
)

type RouteNode struct {
	NodeName string
	Type     int
	Children []*RouteNode
	Callback HTTPListener
}

var routeGetRootMap = make(map[string]*RouteNode)
var routePostRootMap = make(map[string]*RouteNode)

//RegistGetRouter add get func
func RegistGetRouter(url string, f HTTPListener) {
	values := strings.Split(url, "/")
	addNode(nil, values, f, routeGetRootMap)
}

//RegistPostRouter add post func
func RegistPostRouter(url string, f HTTPListener) {
	values := strings.Split(url, "/")
	addNode(nil, values, f, routePostRootMap)
}

func addNode(parent *RouteNode, values []string, f HTTPListener, m map[string]*RouteNode) {
	value := values[0]
	node := new(RouteNode)

	//Type
	if value[0] == ':' {
		node.Type = NODE_TYPE_PARAM
		var name = value[1:len(value)]
		node.NodeName = name
	} else {
		node.Type = NODE_TYPE_COMMAND
		node.NodeName = value
	}

	if parent != nil {
		if parent.Children == nil {
			parent.Children = make([]*RouteNode, 0)
		}
		parent.Children = append(parent.Children, node)
	} else {
		m[value] = node
	}

	if len(values) == 1 {
		node.Callback = f
	} else {
		addNode(node, values[1:len(values)], f, m)
	}
}

//MatchGetURL get url
func MatchGetURL(url string) (HTTPListener, map[string]string, bool) {
	node, m, result := matchURL(url, routeGetRootMap)
	return node.Callback, m, result
}

//MatchPostURL post url
func MatchPostURL(url string) (HTTPListener, map[string]string, bool) {
	node, m, result := matchURL(url, routePostRootMap)
	return node.Callback, m, result
}

//MatchURL url match
func matchURL(url string, m map[string]*RouteNode) (*RouteNode, map[string]string, bool) {
	values := strings.Split(url, "/")
	//get Root
	root, rootOk := m[values[0]]
	if !rootOk {
		return nil, nil, false
	}

	matchList := make([]*MatchCandidate, 0)
	var searchResult bool
	match := make(map[string]string)

	matchList, searchResult = searchFunc(values, root, matchList, match)
	if len(matchList) == 0 && !searchResult {
		return nil, nil, false
	} else if searchResult {
		return root, match, true
	}

	for {
		if len(matchList) == 0 {
			break
		}

		candidate := matchList[len(matchList)-1]
		var ok bool
		matchList, ok = searchFunc(candidate.values,
			candidate.node,
			matchList[0:len(matchList)-1],
			match)

		if ok {
			return candidate.node, match, true
		}
	}

	return nil, nil, false
}

//MatchCandidate match
type MatchCandidate struct {
	values []string
	node   *RouteNode
}

func searchFunc(values []string,
	current *RouteNode,
	matchList []*MatchCandidate,
	match map[string]string) ([]*MatchCandidate, bool) {

	value := values[0]
	end := false

	if len(values) == 1 {
		end = true
	}

	//check param
	if end && len(current.Children) != 0 {
		return nil, false
	}

	if !end && len(current.Children) == 0 {
		return nil, false
	}

	if current.Type == NODE_TYPE_PARAM {
		match[current.NodeName] = value
		if end {
			//found
			return nil, true
		}
	} else if current.Type == NODE_TYPE_COMMAND {
		if current.NodeName == value {
			if end {
				return nil, true
			}
		} else {
			return nil, false
		}
	}

	if current.Children != nil {
		for _, mNode := range current.Children {
			candidate := new(MatchCandidate)
			candidate.values = values[1:len(values)]
			candidate.node = mNode
			matchList = append(matchList, candidate)
		}
	}

	return matchList, false
}
