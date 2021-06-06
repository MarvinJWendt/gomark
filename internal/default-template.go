package internal

import (
	_ "embed"
	"go/doc"
)

type Pkg doc.Package

//go:embed default.tmpl.md
var DefaultMarkdownTemplate string
