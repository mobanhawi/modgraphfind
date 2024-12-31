package main

import (
	"bytes"
	"io"
	"testing"
)

func Test_modgraphfind(t *testing.T) {
	type args struct {
		nodes []string
		in    io.Reader
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
		wantErr bool
	}{
		{
			name: "bad graph - error",
			args: args{
				nodes: []string{"test.com/D@v1.0.0"},
				in: bytes.NewBuffer([]byte(`
test.com/A@v1.0.0 test.com/B@v1.0.0 test.com/B@v1.0.0
test.com/A@v1.0.0 test.com/C@v1.0.0
test.com/B@v1.0.0 test.com/D@v1.0.0
test.com/C@v1.0.0 test.com/E@v1.0.0
`)),
			},
			wantOut: "",
			wantErr: true,
		},
		{
			name: "depth 2 - success",
			args: args{
				nodes: []string{"test.com/D@v1.0.0"},
				in: bytes.NewBuffer([]byte(`
test.com/A@v1.0.0 test.com/B@v1.0.0
test.com/A@v1.0.0 test.com/C@v1.0.0
test.com/B@v1.0.0 test.com/D@v1.0.0
test.com/C@v1.0.0 test.com/E@v1.0.0
`)),
			},
			wantOut: `digraph gomodgraph {
	node [ shape=rectangle fontsize=12 ]
	"test.com/B@v1.0.0" -> "test.com/A@v1.0.0";
	"test.com/D@v1.0.0" -> "test.com/B@v1.0.0";
}
`,
			wantErr: false,
		},
		{
			name: "depth 1 - success",
			args: args{
				nodes: []string{"test.com/C@v1.0.0"},
				in: bytes.NewBuffer([]byte(`
test.com/A@v1.0.0 test.com/B@v1.0.0
test.com/A@v1.0.0 test.com/C@v1.0.0
test.com/B@v1.0.0 test.com/D@v1.0.0
test.com/C@v1.0.0 test.com/E@v1.0.0
`)),
			},
			wantOut: `digraph gomodgraph {
	node [ shape=rectangle fontsize=12 ]
	"test.com/C@v1.0.0" -> "test.com/A@v1.0.0";
}
`,
			wantErr: false,
		},
		{
			name: "depth 1/2 - success",
			args: args{
				nodes: []string{"test.com/B@v1.0.0", "test.com/E@v1.0.0"},
				in: bytes.NewBuffer([]byte(`
test.com/A@v1.0.0 test.com/B@v1.0.0
test.com/A@v1.0.0 test.com/C@v1.0.0
test.com/B@v1.0.0 test.com/D@v1.0.0
test.com/C@v1.0.0 test.com/E@v1.0.0
`)),
			},
			wantOut: `digraph gomodgraph {
	node [ shape=rectangle fontsize=12 ]
	"test.com/B@v1.0.0" -> "test.com/A@v1.0.0";
	"test.com/C@v1.0.0" -> "test.com/A@v1.0.0";
	"test.com/E@v1.0.0" -> "test.com/C@v1.0.0";
}
`,
			wantErr: false,
		},
		{
			name: "depth 1/2 - same path - success",
			args: args{
				nodes: []string{"test.com/B@v1.0.0", "test.com/D@v1.0.0"},
				in: bytes.NewBuffer([]byte(`
test.com/A@v1.0.0 test.com/B@v1.0.0
test.com/A@v1.0.0 test.com/C@v1.0.0
test.com/B@v1.0.0 test.com/D@v1.0.0
test.com/C@v1.0.0 test.com/E@v1.0.0
`)),
			},
			wantOut: `digraph gomodgraph {
	node [ shape=rectangle fontsize=12 ]
	"test.com/B@v1.0.0" -> "test.com/A@v1.0.0";
	"test.com/D@v1.0.0" -> "test.com/B@v1.0.0";
}
`,
			wantErr: false,
		},
		{
			name: "empty nodes - success",
			args: args{
				nodes: []string{""},
				in: bytes.NewBuffer([]byte(`
test.com/A@v1.0.0 test.com/B@v1.0.0
test.com/A@v1.0.0 test.com/C@v1.0.0
test.com/B@v1.0.0 test.com/D@v1.0.0
test.com/C@v1.0.0 test.com/E@v1.0.0
`)),
			},
			wantOut: `digraph gomodgraph {
	node [ shape=rectangle fontsize=12 ]
}
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			err := modgraphfind(tt.args.nodes, tt.args.in, out)
			if (err != nil) != tt.wantErr {
				t.Errorf("modgraphfind() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("modgraphfind() gotOut = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
