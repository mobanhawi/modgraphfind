# ModGraphFind
ModGraphFind parses and searches “go mod graph” output into Graphviz's DOT 


Usage:

    go mod graph | modgraphfind <nodes> | dot -Tpng -o x.png
    go mod graph | modgraphfind <nodes>
