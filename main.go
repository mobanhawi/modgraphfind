// Modgraphfind converts “go mod graph” output into Graphviz's DOT language,
// for use with Graphviz visualization
// Usage:
//
//	go mod graph | modgraphfind <node> | dot -Tpng -o x.png
//	go mod graph | modgraphfind <node>
//
// See also golang.org/x/tools/cmd/digraph for analysis
// of “go mod graph” output.
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: go mod graph | modgraphfind <nodes>`)
	os.Exit(2)
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("modgraphfind: ")

	flag.Usage = usage
	flag.Parse()
	if flag.NArg() == 0 {
		usage()
	}

	if err := modgraphfind(flag.Args(), os.Stdin, os.Stdout); err != nil {
		log.Fatal(err)
	}
}

func modgraphfind(nodes []string, in io.Reader, out io.Writer) error {
	// parse `go mod graph` output into a directed graph
	g, err := parse(in)
	if err != nil {
		return fmt.Errorf("parsing graph: %w", err)
	}

	// add search nodes into a nodeset
	n := nodeset{}
	for _, node := range nodes {
		n[node] = true
	}

	// find subgraph from search nodes to root node
	// TODO: add reverse mode to get all reachable
	// nodes from root node
	sub := g.reachableFrom(n)
	// output subgraph to writer
	var b bytes.Buffer

	err = sub.toDot(&b)
	if err != nil {
		return fmt.Errorf("building graph dot: %w", err)
	}
	_, err = out.Write(b.Bytes())
	if err != nil {
		return fmt.Errorf("writing to stdout: %w", err)
	}
	return nil
}

// parse reads “go mod graph” output from r and returns a graph
func parse(r io.Reader) (graph, error) {
	scanner := bufio.NewScanner(r)

	g := make(graph)

	for scanner.Scan() {
		l := scanner.Text()
		if l == "" {
			continue
		}
		parts := strings.Fields(l)
		if len(parts) != 2 {
			return nil, fmt.Errorf("expected 2 words in line, but got %d: %s", len(parts), l)
		}
		to := parts[0]
		from := parts[1]

		g.addEdges(from, to)
	}
	return g, nil
}
