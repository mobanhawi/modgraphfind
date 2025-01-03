# ModGraphFind
ModGraphFind parses and searches “go mod graph” output into Graphviz's DOT 

![ci/cd]("https://github.com/mobanhawi/modgraphfind/actions/workflows/go.yml/badge.svg") 

![coverage](./assets/coverage.svg)


Inspired by [digraph](golang.org/x/tools/cmd/digraph)


Usage:

    go mod graph | modgraphfind <nodes> | dot -Tpng -o x.png
    go mod graph | modgraphfind <nodes>
