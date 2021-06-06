{{- /*gotype: github.com/MarvinJWendt/gomark/internal.Package*/ -}}
# {{.Name}}

{{if .Doc}}{{.Doc}}{{end}}

{{if or (gt (len .Constants) 0) (gt (len .ConstantBlocks) 0) -}}
## Constants
{{range .Constants}}
### {{.Name}}

```go
{{.Definition}}
```

{{.Doc}}
{{end}}
{{if gt (len .ConstantBlocks) 0}}
## Constant Blocks

{{range .ConstantBlocks -}}
{{if ne .Doc "" }}{{.Doc}}
{{end}}
```go
{{- range .Variables}}
{{.Definition -}}
{{end}}
```

{{end}}{{end}}{{end}}

{{if or (gt (len .Variables) 0) (gt (len .VariableBlocks) 0) -}}
## Variables
{{ range .Variables}}
### {{.Name}}

```go
{{.Definition}}
```

{{.Doc}}
{{end}}
{{if gt (len .VariableBlocks) 0}}
## Variable Blocks

{{range .VariableBlocks -}}
{{if ne .Doc "" }}{{.Doc}}
{{end}}
```go
{{- range .Variables}}
{{.Definition -}}
{{end}}
```

{{end}}{{end}}{{end}}

{{if gt (len .Functions) 0 -}}
## Functions
{{range .Functions}}
### {{.Name}}

```go
{{.Definition}}
```

{{.Doc}}
{{end}}{{end}}

{{if gt (len .Types) 0 -}}
## Types
{{range .Types}}{{$name := .Name}}
### {{.Name}}

```go
{{.Definition}}
```

{{.Doc}}

{{if gt (len .Functions) 0 -}}
{{range .Functions}}
#### {{$name}}.{{.Name}}

```go
{{.Definition}}
```

{{.Doc}}
{{end}}{{end}}{{end}}{{end}}

{{if gt (len .Structs) 0 -}}
## Structs
{{range .Structs}}{{$name := .Name}}
### {{.Name}}

```go
{{.Definition}}
```

{{.Doc}}

{{if gt (len .Functions) 0 -}}
{{range .Functions}}
#### {{$name}}.{{.Name}}

```go
{{.Definition}}
```

{{.Doc}}
{{end}}{{end}}{{end}}{{end}}

{{if gt (len .Interfaces) 0 -}}
## Interfaces
{{range .Interfaces}}
### {{.Name}}

```go
{{.Definition}}
```

{{.Doc}}
{{end}}
{{end}}
