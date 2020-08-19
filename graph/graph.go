package graph

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
	"text/template"
)

var graphTemplate = `digraph {
{{- if eq .direction "horizontal" -}}
rankdir=LR;
{{- end }}
node [shape=box];
{{- $hidePath := .hidePath }}
{{- $hideVersion := .hideVersion }}
{{- range $modId, $mod := .mods }}
{{ $mod.ID }} [label="{{ if not $hidePath }}{{ $mod.Path }}/{{ end -}}{{ $mod.Name }}{{ if not $hideVersion }}{{printf "\n"}}{{ $mod.Version }}{{ end -}}"]
{{- end }}
{{ range $modId, $depModIds := .dependencies -}}
{{- range $_, $depModId := $depModIds -}}
{{ $modId }} -> {{ $depModId }};
{{  end -}}
{{- end -}}
}
`

// ModuleGraph represents a graph of module dependencies
type ModuleGraph struct {
	Reader io.Reader

	Mods         map[string]*module
	Dependencies map[int][]int
}

type module struct {
	ID      int
	Path    string
	Name    string
	Version string
}

type modulesExclusionArray []string

func (a *modulesExclusionArray) String() string {
	if a != nil {
		return strings.Join(*a, ",")
	}
	return ""
}

func (a *modulesExclusionArray) Set(s string) error {
	*a = strings.Split(strings.ReplaceAll(s, " ", ""), ",")
	return nil
}

type ParseOptions struct {
	IgnoreVersion bool // if true, all versions are seen as one
	IgnoreModules modulesExclusionArray
	IgnoreIndirect bool
}

type RenderOptions struct {
	HidePath    bool // if true, does not show module's path
	HideVersion bool
}

// NewModuleGraph creates a new ModuleGraph
func NewModuleGraph(r io.Reader) *ModuleGraph {
	return &ModuleGraph{
		Reader: r,

		Mods:         make(map[(string)]*module),
		Dependencies: make(map[int][]int),
	}
}

func stripModulePath(mod string) string {
	index := strings.LastIndexAny(mod, "/")
	if index == -1 {
		return mod
	}
	return mod[index+1 : len(mod)]
}

func stripModuleVersion(mod string) string {
	return strings.TrimRight(strings.SplitAfter(mod, "@")[0], "@")
}

func getModulePath(mod string) string {
	mod = stripModuleVersion(mod)

	index := strings.LastIndexAny(mod, "/")
	if index == -1 {
		return ""
	}
	return mod[0:index]
}

func getModuleVersion(mod string) string {
	index := strings.LastIndexAny(mod, "@")
	if index == -1 {
		return ""
	}
	return mod[index+1 : len(mod)]
}

func getModuleName(mod string) string {
	return stripModulePath(stripModuleVersion(mod))
}

func (m *ModuleGraph) hasDependency(from, to int) bool {
	count := len(m.Dependencies[from])
	for i := 0; i < count; i++ {
		if m.Dependencies[from][i] == to {
			return true
		}
	}
	return false
}

func isStringInSortedArray(s string, a []string) bool {
	if !sort.StringsAreSorted(a) {
		panic("Strings expected to be sorted")
	}

	count := len(a)
	index := sort.Search(count, func(i int) bool { return a[i] >= s })
	return index < count && a[index] == s
}

func (m *ModuleGraph) hasModule(mod string) bool {
	_, ok := m.Mods[mod]
	return ok
}

// Parse parses the input into a ModuleGraph
func (m *ModuleGraph) Parse(options ParseOptions) error {
	bufReader := bufio.NewReader(m.Reader)

	var topModule string

	serialID := 1
	for {
		relationBytes, err := bufReader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return err
		}

		if len(relationBytes) <= 0 {
			return nil
		}

		relation := bytes.Split(relationBytes, []byte(" "))
		mod, depMod := strings.TrimSpace(string(relation[0])), strings.TrimSpace(string(relation[1]))

		modifiedMod := mod
		modifiedDepMod := depMod

		if options.IgnoreVersion {
			modifiedMod = stripModuleVersion(mod)
			modifiedDepMod = stripModuleVersion(depMod)
		}

		if topModule == "" {
			topModule = modifiedMod
		}

		isIndirect := topModule != modifiedMod
		if options.IgnoreIndirect && isIndirect {
			continue
		}

		var md, depmd *module
		var ok bool

		if !isStringInSortedArray(modifiedMod, options.IgnoreModules) {
			md, ok = m.Mods[modifiedMod]
			if !ok {
				md = &module{ID: serialID, Path: getModulePath(mod), Name: getModuleName(mod), Version: getModuleVersion(mod)}
				m.Mods[modifiedMod] = md
				serialID++
			}
		}

		if !isStringInSortedArray(modifiedDepMod, options.IgnoreModules) {
			depmd, ok = m.Mods[modifiedDepMod]
			if !ok {
				depmd = &module{ID: serialID, Path: getModulePath(depMod), Name: getModuleName(depMod), Version: getModuleVersion(depMod)}
				m.Mods[modifiedDepMod] = depmd
				serialID++
			}
		}

		if md != nil && depmd != nil && !m.hasDependency(md.ID, depmd.ID) {
			m.Dependencies[md.ID] = append(m.Dependencies[md.ID], depmd.ID)
		}
	}
}

// Render renders ModuleGraph into dot tool syntax
func (m *ModuleGraph) Render(w io.Writer, options RenderOptions) error {
	templ, err := template.New("graph").Parse(graphTemplate)
	if err != nil {
		return fmt.Errorf("templ.Parse: %v", err)
	}

	var direction string
	if len(m.Dependencies) > 15 {
		direction = "horizontal"
	}

	if err := templ.Execute(w, map[string]interface{}{
		"mods":         m.Mods,
		"dependencies": m.Dependencies,
		"direction":    direction,
		"hidePath":     options.HidePath,
		"hideVersion":  options.HideVersion,
	}); err != nil {
		return fmt.Errorf("templ.Execute: %v", err)
	}

	return nil
}
