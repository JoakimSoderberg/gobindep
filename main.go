package main

import (
	"github.com/mitchellh/golicense/module"
	"github.com/rsc/goversion/version"

	"fmt"
	"os"

	flag "github.com/spf13/pflag"
)

func main() {
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	flags.Parse(os.Args[1:])
	args := flags.Args()

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

	for _, mod := range mods {
		fmt.Printf("%s %s\n", mod.Path, mod.Version)
	}
}
