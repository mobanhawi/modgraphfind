# ModGraphFind
ModGraphFind parses and searches “go mod graph” output into Graphviz's DOT 

Inpsired by [digraph](golang.org/x/tools/cmd/digraph)


Usage:

    go mod graph | modgraphfind <nodes> | dot -Tpng -o x.png
    go mod graph | modgraphfind <nodes>
