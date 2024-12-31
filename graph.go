// See also golang.org/x/tools/cmd/digraph for general queries and analysis
// of “go mod graph” output.

package main

import (
	"bytes"
	"fmt"
	"sort"
)

// node represents a node in a graph
type node string

// newNode creates a new node
func newNode(s string) node {
	return node(s)
}

// nodelist represents nodes as strings
type nodelist []node

// nodeset represents nodes as strings
type nodeset map[node]bool

// graph maps nodes to the non-nil set of their immediate successors.
type graph map[node]nodeset

func (s nodeset) sort() nodelist {
	nodes := make(nodelist, len(s))
	var i int
	for n := range s {
		nodes[i] = n
		i++
	}
	// ensure deterministic order
	// we are looping over nodeset which is a map
	sort.SliceStable(nodes, func(i, j int) bool {
		return nodes[i] < nodes[j]
	})
	return nodes
}

func (g graph) addNode(node node) nodeset {
	edges := g[node]
	if edges == nil {
		edges = make(nodeset)
		g[node] = edges
	}
	return edges
}

func (g graph) addEdges(from node, to ...node) {
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
	var visit func(node)
	visit = func(n node) {
		if !seen[n] {
			seen[n] = true
			// visit all immediate successors
			for e := range g[n] {
				sub.addEdges(n, e)
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
	fmt.Fprintln(w, "digraph gomodgraph {")
	fmt.Fprintln(w, "\tnode [ shape=rectangle fontsize=12 ]")

	for _, src := range g.nodelist() {
		for _, dst := range g[src].sort() {
			_, err = fmt.Fprintf(w, "\t%q -> %q;\n", src, dst)
			if err != nil {
				return
			}
		}
	}
	fmt.Fprintln(w, "}")
	return
}
