package GoaTS

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/goadesign/goa/design"
	"github.com/goadesign/goa/goagen/codegen"
)

//Generator is the code generator type
type Generator struct {
	genfiles   []string // Generated files
	outDir     string   // Absolute path to output directory
	target     string   // Target package name - "models" by default
	appPkg     string   // Generated goa app package name - "app" by default
	appPkgPath string   // Generated goa app package import path
}

//Generate is the generator entry point called by the meta generator.
//Errors found in the program should probably be passed up the tree and let
//goagen evaluate it
func Generate() (files []string, err error) {
	var outDir, target, appPkg, ver string

	set := flag.NewFlagSet("goats", flag.PanicOnError)
	set.String("design", "", "")
	set.StringVar(&outDir, "out", "", "")
	set.StringVar(&ver, "version", "", "")
	set.StringVar(&target, "pkg", "interfaces", "")
	set.StringVar(&appPkg, "app", "app", "")
	set.Parse(os.Args[2:])

	// First check compatibility
	if err := codegen.CheckVersion(ver); err != nil {
		return nil, err
	}

	// Now proceed
	appPkgPath, err := codegen.PackagePath(filepath.Join(outDir, appPkg))
	if err != nil {
		return nil, fmt.Errorf("GoaTS: invalid app package: %s", err)
	}

	g := &Generator{outDir: outDir, target: target, appPkg: appPkg, appPkgPath: appPkgPath}

	//generate a data structure describing each of the resources including their payloads
	r, err := g.genResourceNames(design.Design)
	if err != nil {
		return []string{}, err
	}

	//iterate over the resource structure
	for _, a := range r {
		//for each resource, write the payloads to a file
		f, err := g.writeTypes(a.Name, a.Payloads)
		if err != nil {
			return []string{}, nil
		}
		if len(f) > 0 {
			fmt.Printf("Created file of typescript interfaces for %s resource: %s\n", a.Name, f)
		}
	}

	return []string{}, nil
}
