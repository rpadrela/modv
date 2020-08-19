package graph

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"sort"
	"strings"
	"testing"
)

// to remove
func assertStringInArray(t *testing.T, s string, a []string) {
	if !isStringInSortedArray(s, a) {
		t.Error("string ", s, " expected to be in array")
	}
}

// to remove
func assertStringNotInArray(t *testing.T, s string, a []string) {
	if isStringInSortedArray(s, a) {
		t.Error("string ", s, " not expected to be in array")
	}
}

func TestGetModuleName(t *testing.T) {
	assert.Equal(t, getModuleName("golang.org/x/sys@v0.0.0-20191010194322-b09"), "sys", "Unexpected module name")
	assert.Equal(t, getModuleName("github.com/fatih/color@v1.7.0"), "color", "Unexpected module name")
}

func TestGetModulePath(t *testing.T) {
	assert.Equal(t, getModulePath("golang.org/x/sys@v0.0.0-20191010194322-b09"), "golang.org/x", "Unexpected module path")
	assert.Equal(t, getModulePath("github.com/fatih/color@v1.7.0"), "github.com/fatih", "Unexpected module path")
}

func TestGetModuleVersion(t *testing.T) {
	assert.Equal(t, getModuleVersion("golang.org/x/sys@v0.0.0-20191010194322-b09"), "v0.0.0-20191010194322-b09", "Unexpected module version")
	assert.Equal(t, getModuleVersion("github.com/fatih/color@v1.7.0"), "v1.7.0", "Unexpected module version")
}

// Consider renaming test name?
func TestIsStringInSortedArray(t *testing.T) {
	a := []string{
		"World",
		"Hello",
		"Hello World",
		"One",
		"Two",
		"Three",
	}

	sort.Strings(a)

	assertStringInArray(t, "World", a)
	assertStringInArray(t, "Three", a)
	assertStringNotInArray(t, "world", a)
	assertStringNotInArray(t, "three", a)
	assertStringNotInArray(t, "nope", a)
}

func TestGraphContainsAllModules(t *testing.T) {
	dependencies := []string {
		"github.com/poloxue/testmod golang.org/x/text@v0.3.2",
		"github.com/poloxue/testmod rsc.io/quote/v3@v3.1.0",
		"github.com/poloxue/testmod rsc.io/sampler@v1.3.1",
		"golang.org/x/text@v0.3.2 golang.org/x/tools@v0.0.0-20180917221912-90fa682c2a6e",
		"rsc.io/quote/v3@v3.1.0 rsc.io/sampler@v1.3.0",
		"rsc.io/sampler@v1.3.1 golang.org/x/text@v0.0.0-20170915032832-14c0d48ead0c",
		"rsc.io/sampler@v1.3.0 golang.org/x/text@v0.0.0-20170915032832-14c0d48ead0c",
	}

	expectedModules := []string{
		"github.com/poloxue/testmod",
		"golang.org/x/text@v0.3.2",
		"rsc.io/quote/v3@v3.1.0",
		"rsc.io/sampler@v1.3.0",
		"rsc.io/sampler@v1.3.1",
		"golang.org/x/text@v0.0.0-20170915032832-14c0d48ead0c",
		"golang.org/x/tools@v0.0.0-20180917221912-90fa682c2a6e",
	}

	optionsIgnoringVersion := ParseOptions{IgnoreVersion: false}

	moduleGraph := NewModuleGraph(bytes.NewReader([]byte(strings.Join(dependencies[:],"\n"))))
	moduleGraph.Parse(optionsIgnoringVersion)

	assert.Equal(t, len(expectedModules), len(moduleGraph.Mods), "Number of expected modules is incorrect")

	for _, expectedModule := range expectedModules {
		assert.True(t, moduleGraph.hasModule(expectedModule), "Module \"" + expectedModule + "\" expected in graph")
	}
}

func TestGraphContainsAllModulesIgnoreVersion(t *testing.T) {
	dependencies := []string {
		"github.com/poloxue/testmod golang.org/x/text@v0.3.2",
		"github.com/poloxue/testmod rsc.io/quote/v3@v3.1.0",
		"github.com/poloxue/testmod rsc.io/sampler@v1.3.1",
		"golang.org/x/text@v0.3.2 golang.org/x/tools@v0.0.0-20180917221912-90fa682c2a6e",
		"rsc.io/quote/v3@v3.1.0 rsc.io/sampler@v1.3.0",
		"rsc.io/sampler@v1.3.1 golang.org/x/text@v0.0.0-20170915032832-14c0d48ead0c",
		"rsc.io/sampler@v1.3.0 golang.org/x/text@v0.0.0-20170915032832-14c0d48ead0c",
	}

	expectedModules := []string{
		"github.com/poloxue/testmod",
		"golang.org/x/text@v0.3.2",
		"rsc.io/quote/v3@v3.1.0",
		"rsc.io/sampler@v1.3.0",
		"rsc.io/sampler@v1.3.1",
		"golang.org/x/text@v0.0.0-20170915032832-14c0d48ead0c",
		"golang.org/x/tools@v0.0.0-20180917221912-90fa682c2a6e",
	}

	optionsNotIgnoringVersion := ParseOptions{IgnoreVersion: false}

	moduleGraph := NewModuleGraph(bytes.NewReader([]byte(strings.Join(dependencies[:],"\n"))))
	moduleGraph.Parse(optionsNotIgnoringVersion)

	assert.Equal(t, len(expectedModules), len(moduleGraph.Mods), "Number of expected modules is incorrect")

	for _, expectedModule := range expectedModules {
		assert.True(t, moduleGraph.hasModule(expectedModule), "Module \"" + expectedModule + "\" expected in graph")
	}
}


/*
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "ignoring version",
			args: args{
				bytes.NewReader([]byte(`github.com/poloxue/testmod golang.org/x/text@v0.3.2
github.com/poloxue/testmod rsc.io/quote/v3@v3.1.0
github.com/poloxue/testmod rsc.io/sampler@v1.3.1
golang.org/x/text@v0.3.2 golang.org/x/tools@v0.0.0-20180917221912-90fa682c2a6e
rsc.io/quote/v3@v3.1.0 rsc.io/sampler@v1.3.0
rsc.io/sampler@v1.3.1 golang.org/x/text@v0.0.0-20170915032832-14c0d48ead0c
rsc.io/sampler@v1.3.0 golang.org/x/text@v0.0.0-20170915032832-14c0d48ead0c`))},
			want: []string{
				"github.com/poloxue/testmod",
				"golang.org/x/text",
				"rsc.io/quote/v3",
				"rsc.io/sampler",
				"golang.org/x/tools",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			moduleGraph := NewModuleGraph(tt.args.reader)
			moduleGraph.Parse(ParseOptions{IgnoreVersion: true })
			assertHasExactModules(t, moduleGraph, tt.want)
		})
	}
*/