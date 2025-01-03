# ModGraphFind
ModGraphFind parses and searches “go mod graph” output into Graphviz's DOT 

<p>
    <a href="https://github.com/mobanhawi/modgraphfind/actions/workflows/go.yml"><img src="https://github.com/mobanhawi/modgraphfind/actions/workflows/go.yml/badge.svg"></a>
    <a href="https://github.com/mobanhawi/modgraphfind/actions/workflows/go.yml"><img src="https://github.com/mobanhawi/modgraphfind/wiki/coverage.svg"></a>
</p>

Inspired by [digraph](golang.org/x/tools/cmd/digraph)


Usage:

    go mod graph | modgraphfind <nodes> | dot -Tpng -o x.png
    go mod graph | modgraphfind <nodes>
