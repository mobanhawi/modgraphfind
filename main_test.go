package main

import (
	"bytes"
	"testing"
)

func TestRun(t *testing.T) {
	out := &bytes.Buffer{}
	in := bytes.NewBuffer([]byte(`
test.com/A@v1.0.0 test.com/B@v1.0.0
test.com/A@v1.0.0 test.com/C@v1.0.0
test.com/B@v1.0.0 test.com/D@v1.0.0
test.com/C@v1.0.0 test.com/E@v1.0.0
`))
	if err := modgraphfind([]string{"test.com/D@v1.0.0"}, in, out); err != nil {
		t.Fatal(err)
	}

	gotGraph := out.String()
	wantGraph := `digraph gomodgraph {
	node [ shape=rectangle fontsize=12 ]
	"test.com/B@v1.0.0" -> "test.com/A@v1.0.0";
	"test.com/D@v1.0.0" -> "test.com/B@v1.0.0";
}
`
	if gotGraph != wantGraph {
		t.Fatalf("\ngot: %s\nwant: %s", gotGraph, wantGraph)
	}
}
