// Command go-bootstrap scaffolds new Go projects.
package main

import (
	"embed"
	"io/fs"
	"os"

	"github.com/npikall/go-bootstrap/internal/cli"
)

var version = "dev"

//go:embed all:template
var templateFS embed.FS

func main() {
	tmpl, err := fs.Sub(templateFS, "template")
	if err != nil {
		panic(err)
	}
	os.Exit(cli.Run(os.Args[1:], tmpl, version, os.Stdout, os.Stderr))
}
