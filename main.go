package main

import (
	"encoding/json"
	"io/ioutil"
	"text/template"

	"github.com/JoakimSoderberg/gobindep/module"
	"github.com/rsc/goversion/version"

	"fmt"
	"os"

	flag "github.com/spf13/pflag"
)

// Output represents an output context.
//
// This is used when outputting
type Output struct {
	Executable string          `json:"executable"`
	Size       int64           `json:"size"`
	Modules    []module.Module `json:"modules"`
}

var replaceOverwrites bool

func main() {
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	showHelp := flags.Bool("help", false, "Show this help")
	outputJSON := flags.Bool("json", false, "Output results in JSON")
	outputTemplate := flags.String("template", "", "Output results based on the specified template string")
	outputTemplateFile := flags.String("template-file", "", "Output results based on the specified template file")
	flags.BoolVar(&replaceOverwrites, "replace-overwrites", true, "When outputting should a replaced module overwrite the original? If set to false, we add a nested \"replace\" instead for JSON.")

	flags.Parse(os.Args[1:])
	args := flags.Args()

	if *showHelp {
		flags.PrintDefaults()
		os.Exit(0)
	}

	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Please provide an go executable path\n")
		os.Exit(1)
	}

	exePath := args[0]

	// Read the dependencies from the binary itself
	vsn, err := version.ReadExe(exePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %q: %s\n", args[0], err)
		os.Exit(1)
	}

	if vsn.ModuleInfo == "" {
		fmt.Fprintf(os.Stderr, "This executable was either compiled without Go modules or has no dependencies")
		os.Exit(1)
	}

	// From the raw module string from the binary, we need to parse this
	// into structured data with the module information.
	mods, err := module.ParseExeData(vsn.ModuleInfo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing dependencies: %s\n", err)
		os.Exit(1)
	}

	if replaceOverwrites {
		for i := range mods {
			mod := &mods[i]
			if mod.Replace == nil {
				continue
			}

			mod.Path = mod.Replace.Path
			mod.Version = mod.Replace.Version
			mod.Hash = mod.Replace.Hash
			mod.Replace = nil
		}
	}

	var output Output
	output.Executable = exePath
	output.Modules = mods

	// Read the size of the executable.
	fileinfo, err := os.Stat(exePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error stating file: %s\n", err)
		os.Exit(1)
	}

	output.Size = fileinfo.Size()

	var templateString string
	if len(*outputTemplateFile) > 0 {
		contents, err := ioutil.ReadFile(*outputTemplateFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error, failed to open output template: %s\n", err)
			os.Exit(1)
		}
		templateString = string(contents)
	} else {
		templateString = *outputTemplate
	}

	if len(templateString) > 0 {
		tmpl := template.New("template")

		tmpl, err = tmpl.Parse(string(templateString))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error, failed to parse the output template: %s\n", err)
			os.Exit(1)
		}

		err = tmpl.Execute(os.Stdout, output)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error, failed to execute template: %s\n", err)
			os.Exit(1)
		}
	} else if *outputJSON {
		json, err := json.Marshal(output)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to output module info as JSON: %s\n", err)
			os.Exit(1)
		}

		fmt.Println(string(json))
	} else {
		for _, mod := range mods {
			fmt.Printf("%s %s\n", mod.Path, mod.Version)
		}
	}
}
