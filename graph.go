// See also golang.org/x/tools/cmd/digraph for general queries and analysis
// of “go mod graph” output.

package main

import (
	"bytes"
	"fmt"
	"sort"
)

// nodelist represents nodes as strings
type nodelist []string

// nodeset represents nodes as strings
type nodeset map[string]bool

// graph maps nodes to the non-nil set of their immediate successors.
type graph map[string]nodeset

func (s nodeset) sort() nodelist {
	nodes := make(nodelist, len(s))
	var i int
	for node := range s {
		nodes[i] = node
		i++
	}
	sort.Strings(nodes)
	return nodes
}

func (g graph) addNode(node string) nodeset {
	edges := g[node]
	if edges == nil {
		edges = make(nodeset)
		g[node] = edges
	}
	return edges
}

func (g graph) addEdges(from string, to ...string) {
	edges := g.addNode(from)
	for _, to := range to {
		g.addNode(to)
		edges[to] = true
	}
}

func (g graph) nodelist() nodelist {
	nodes := make(nodeset)
	for node := range g {
		nodes[node] = true
	}
	return nodes.sort()
}

// reachableFrom subgraph transitively reach the specified nodes
func (g graph) reachableFrom(roots nodeset) graph {
	seen := make(nodeset)
	sub := make(graph)
	var visit func(node string)
	visit = func(node string) {
		if !seen[node] {
			seen[node] = true
			// visit all immediate successors
			for e := range g[node] {
				sub.addEdges(node, e)
				visit(e)
			}
		}
	}
	for root := range roots {
		visit(root)
	}
	return sub
}

// transpose the reverse of the input edges
// nolint
func (g graph) transpose() graph {
	rev := make(graph)
	for node, edges := range g {
		rev.addNode(node)
		for succ := range edges {
			rev.addEdges(succ, node)
		}
	}
	return rev
}

// toDot TODO: extend to add colors
func (g graph) toDot(w *bytes.Buffer) (err error) {
	_, err = fmt.Fprintln(w, "digraph gomodgraph {")
	if err != nil {
		return
	}
	_, err = fmt.Fprintln(w, "\tnode [ shape=rectangle fontsize=12 ]")
	if err != nil {
		return
	}

	for _, src := range g.nodelist() {
		for _, dst := range g[src].sort() {
			_, err = fmt.Fprintf(w, "\t%q -> %q;\n", src, dst)
			if err != nil {
				return
			}
		}
	}
	_, err = fmt.Fprintln(w, "}")
	if err != nil {
		return
	}
	return
}
