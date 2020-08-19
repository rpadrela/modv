package main

import (
	"flag"
	"fmt"
	"github.com/rpadrela/modv/graph"
	"os"
	"runtime"
	"sort"
)

func PrintUsage() {
	fmt.Printf("\nUsages:\n\n")
	switch runtime.GOOS {
	case "darwin":
		fmt.Printf("\tgo mod graph | modv | dot -T svg | open -f -a /System/Applications/Preview.app")
	case "linux":
		fmt.Printf("\tgo mod graph | modv | dot -T svg -o /tmp/modv.svg | xdg-open /tmp/modv.svg")
	case "windows":
		fmt.Printf("\tgo mod graph | modv | dot -T png -o graph.png; start graph.png")
	}

	fmt.Printf("\n\n")
}

func parseFlags() (graph.ParseOptions, graph.RenderOptions) {
	parseOptions := graph.ParseOptions{}

	ignoreVersionPtr := flag.Bool("ignoreVersion", false, "if true, all versions of the same module will be treated as one")
	//ignoreIndirectPtr := flag.Bool("ignoreIndirect", false, "if true, graph will only show module's direct dependencies")
	hidePathPtr := flag.Bool("hidePath", false, "if true, graph will not display the path of the module")
	flag.Var(&parseOptions.IgnoreModules, "ignoreModules", "comma-separated list of modules to ignore including path (e.g. github.com/x/sys")
	flag.Parse()

	sort.Strings(parseOptions.IgnoreModules)

	if ignoreVersionPtr != nil {
		parseOptions.IgnoreVersion = *ignoreVersionPtr
	}

	if ignoreIndirectPtr != nil {
		parseOptions.IgnoreIndirect = *ignoreIndirectPtr
	}

	renderOptions := graph.RenderOptions{}
	renderOptions.HideVersion = parseOptions.IgnoreVersion
	if hidePathPtr != nil {
		renderOptions.HidePath = *hidePathPtr
	}

	return parseOptions, renderOptions
}

func main() {
	info, err := os.Stdin.Stat()
	if err != nil {
		fmt.Println("os.Stdin.Stat:", err)
		PrintUsage()
		os.Exit(1)
	}

	if info.Mode()&os.ModeNamedPipe == 0 {
		fmt.Println("command err: command is intended to work with pipes.")
		PrintUsage()
		os.Exit(1)
	}

	parseOptions, renderOptions := parseFlags()

	mg := graph.NewModuleGraph(os.Stdin)
	if err := mg.Parse(parseOptions); err != nil {
		fmt.Println("mg.Parse: ", err)
		PrintUsage()
		os.Exit(1)
	}

	if err := mg.Render(os.Stdout, renderOptions); err != nil {
		fmt.Println("mg.Render: ", err)
		PrintUsage()
		os.Exit(1)
	}
}
