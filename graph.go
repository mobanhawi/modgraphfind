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

// list represents nodes as strings
type list []node

// set represents nodes as strings
type set map[node]bool

// graph maps nodes to the non-nil set of their immediate successors.
type graph map[node]set

func (s set) sort() list {
	nodes := make(list, len(s))
	var i int
	for n := range s {
		nodes[i] = n
		i++
	}
	// ensure deterministic order
	// we are looping over set which is a map
	sort.SliceStable(nodes, func(i, j int) bool {
		return nodes[i] < nodes[j]
	})
	return nodes
}

func (g graph) addNode(node node) set {
	edges := g[node]
	if edges == nil {
		edges = make(set)
		g[node] = edges
	}
	return edges
}

func (g graph) addEdges(from node, to ...node) {
	edges := g.addNode(from)
	for _, t := range to {
		g.addNode(t)
		edges[t] = true
	}
}

func (g graph) nodelist() list {
	nodes := make(set)
	for n := range g {
		nodes[n] = true
	}
	return nodes.sort()
}

// reachableFrom subgraph transitively reach the specified nodes
func (g graph) reachableFrom(roots set) graph {
	seen := make(set)
	// reachable subgraph from roots
	sub := make(graph)
	var visit func(node)
	visit = func(n node) {
		if !seen[n] {
			seen[n] = true
			// visit all immediate successors
			for e := range g[n] {
				// add edge from `e` to `n` to subgraph
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
	for n, edges := range g {
		rev.addNode(n)
		for succ := range edges {
			rev.addEdges(succ, n)
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
