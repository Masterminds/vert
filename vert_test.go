package main

import (
	"bytes"
	"testing"

	"github.com/codegangsta/cli"
)

type mockContext struct {
	bools map[string]bool
	args  cli.Args
}

func (c *mockContext) Args() cli.Args {
	return c.args
}
func (c *mockContext) Bool(name string) bool {
	return c.bools[name]
}

func TestGit2semver(t *testing.T) {
	tests := map[string]string{
		"v1.10.0-123-g0239788": "1.10.0+123.g0239788",
		"1.10.0-123-g0239788":  "1.10.0+123.g0239788",
	}

	for in, expect := range tests {
		out, err := git2semver(in)
		if err != nil {
			t.Errorf("Failed to parse %s: %s", in, err)
			continue
		}
		if out.String() != expect {
			t.Errorf("Expected %q, got %q", expect, out.String())
		}
	}

	// And test a failure
	in := "fatal: No names found, cannot describe anything."
	out, err := git2semver(in)
	if err == nil {
		t.Errorf("Expected version parse to fail for %q", in)
	}
	if out != nil {
		t.Errorf("Expected version to be nil for %q", in)
	}
}

func TestRun(t *testing.T) {
	var b bytes.Buffer

	c := &mockContext{
		args: cli.Args{">=1.0.0", "1.0.0", "1.1.1", "1.2.3", "1.0.1", "0.9.0"},
		bools: map[string]bool{
			"failed": false,
			"sort":   true,
		},
	}

	// Set the package defaults
	stdout = &b
	stderr = &b

	tests := []struct {
		args  cli.Args
		bools map[string]bool
		out   string
		code  int
	}{
		// Base case.
		{
			args:  cli.Args{"v1.0.0", "1.0.0"},
			bools: map[string]bool{"failed": false, "sort": false},
			code:  0,
			out:   "1.0.0\n",
		},
		// One failure, four passes, sorted.
		{
			args:  cli.Args{">=1.0.0", "1.0.0", "1.1.1", "1.2.3", "1.0.1", "0.9.0"},
			bools: map[string]bool{"failed": false, "sort": true},
			code:  1,
			out:   "1.0.0\n1.0.1\n1.1.1\n1.2.3\n",
		},
		// One failure, four passes, unsorted.
		{
			args:  cli.Args{">=1.0.0", "1.0.0", "1.1.1", "1.2.3", "1.0.1", "0.9.0"},
			bools: map[string]bool{"failed": false, "sort": false},
			code:  1,
			out:   "1.0.0\n1.1.1\n1.2.3\n1.0.1\n",
		},
		// One failure, print failures.
		{
			args:  cli.Args{">=1.0.0", "1.0.0", "1.1.1", "1.2.3", "1.0.1", "0.9.0"},
			bools: map[string]bool{"failed": true, "sort": true},
			code:  1,
			out:   "0.9.0\n",
		},
		// Two failures, sorted.
		{
			args:  cli.Args{">=1.0.0", "0.1", "v0.9.0"},
			bools: map[string]bool{"failed": true, "sort": true},
			code:  2,
			out:   "0.1.0\n0.9.0\n",
		},
		// Convert git tag
		{
			args:  cli.Args{">1", "v1.10.0-123-g0239788"},
			bools: map[string]bool{"git": true},
			code:  0,
			out:   "1.10.0+123.g0239788\n",
		},
	}

	for _, tt := range tests {
		c.args = tt.args
		c.bools = tt.bools
		res := run(c)
		if res != tt.code {
			t.Errorf("Expected code %d, got %d", tt.code, res)
		}
		if b.String() != tt.out {
			t.Errorf("Expected:%s\nGot:%s", tt.out, b.String())
		}
		b.Reset()
	}

}
