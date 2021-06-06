package internal

import (
	"fmt"
	"strings"
	"unicode"
)

type GoDoc struct {
	Raw      string
	Sections map[string]string
	Package  Package
}

func (d *GoDoc) Parse() error {
	// Parse sections
	if d.Sections == nil {
		d.Sections = make(map[string]string)
	}
	lines := strings.Split(d.Raw, "\n")
	currentSection := "docs"
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Check if line is uppercase
		if IsUpper(trimmedLine) {
			currentSection = strings.ToLower(trimmedLine)
		} else {
			d.Sections[currentSection] += line + "\n"
		}
	}

	// Parse package docs
	docs := d.Sections["docs"]
	lines = strings.Split(docs, "\n")
	d.Package.Name = strings.Fields(lines[0])[1]
	d.Package.Doc = strings.TrimSpace(strings.Join(lines[2:], "\n"))

	// Parse function docs
	docs = d.Sections["functions"]
	lines = strings.Split(docs, "\n")
	var funcs []Function

	for _, line := range lines {
		if getFunctionName(line) != "" {
			funcs = append(funcs, Function{
				Name:       getFunctionName(line),
				Definition: strings.TrimSpace(line),
			})
		} else if line != "" {
			if len(funcs) > 0 {
				funcs[len(funcs)-1].Doc += strings.TrimLeft(line, " ") + "\n"
			}
		}
	}

	d.Package.Functions = funcs

	// Parse variable and constant docs
	for section, short := range map[string]string{"variables": "var", "constants": "const"} {
		docs := d.Sections[section]
		lines := strings.Split(docs, "\n")
		isBlock := false
		var lastDocumentable documentable
		for _, line := range lines {
			s := strings.TrimSpace(line)

			if strings.HasPrefix(line, "    ") {
				if lastDocumentable != nil {
					lastDocumentable.addToDocs(line)
				}
				continue
			}

			if strings.HasPrefix(line, fmt.Sprintf("%s (", short)) {
				isBlock = true
				if short == "var" {
					d.Package.VariableBlocks = append(d.Package.VariableBlocks, VariableBlock{})
					lastDocumentable = d.Package.getLastVariableBlock()
				} else if short == "const" {
					d.Package.ConstantBlocks = append(d.Package.ConstantBlocks, VariableBlock{})
					lastDocumentable = d.Package.getLastConstantBlock()
				}
				continue
			}

			if strings.HasPrefix(line, ")") {
				isBlock = false
				continue
			}

			if isBlock {
				var block *VariableBlock

				if short == "var" {
					block = d.Package.getLastVariableBlock()
				} else if short == "const" {
					block = d.Package.getLastConstantBlock()
				}

				block.Variables = append(block.Variables, parseVariable(line))
			} else if strings.HasPrefix(s, "var") || strings.HasPrefix(s, "const") {
				if short == "var" {
					d.Package.Variables = append(d.Package.Variables, parseVariable(line))
					lastDocumentable = d.Package.getLastVariable()
				} else if short == "const" {
					d.Package.Constants = append(d.Package.Constants, parseVariable(line))
					lastDocumentable = d.Package.getLastConstant()
				}
			}
		}
	}

	// Parse type docs
	docs = d.Sections["types"]
	lines = strings.Split(docs, "\n")

	currentType := ""
	var lastDefinition *string
	var lastFunc *Function
	var lastType string
	var lastDocumentable documentable

	for _, line := range lines {
		s := strings.TrimSpace(line)
		if line == "}" {
			currentType = ""
			*lastDefinition += "}"
			continue
		}

		switch currentType {
		case "struct":
			lastDefinition = &d.Package.getLastStruct().Definition
			lastDocumentable = d.Package.getLastStruct()
			*lastDefinition += line + "\n"

			lastType = currentType
		case "interface":
			lastDefinition = &d.Package.getLastInterface().Definition
			d.Package.getLastInterface().Values = append(d.Package.getLastInterface().Values, parseVariable(s))
			lastDocumentable = d.Package.getLastInterface()
			*lastDefinition += line + "\n"
		case "var":
			lastDocumentable = d.Package.getLastVariable()

			currentType = ""
		case "func":
			if lastType == "struct" {
				d.Package.getLastStruct().Functions = append(d.Package.getLastStruct().Functions, Function{})
				lastFunc = &d.Package.getLastStruct().Functions[len(d.Package.getLastStruct().Functions)-1]
				if lastDefinition != nil && *lastDefinition != "" {
					lastFunc.Name = getFunctionName(*lastDefinition)
					lastFunc.Definition = *lastDefinition
				}
			} else if lastType == "type" {
				d.Package.getLastType().Functions = append(d.Package.getLastType().Functions, Function{})
				lastFunc = &d.Package.getLastType().Functions[len(d.Package.getLastType().Functions)-1]
				if lastDefinition != nil && *lastDefinition != "" {
					lastFunc.Name = getFunctionName(*lastDefinition)
					lastFunc.Definition = *lastDefinition
				}
			}

			lastDocumentable = lastFunc
			currentType = ""
		case "type":
			lastDocumentable = d.Package.getLastType()
			lastType = currentType
			currentType = ""
		case "docs":
			currentType = ""
		}

		switch {
		case strings.HasPrefix(line, "type"):
			switch {
			case strings.Contains(line, "struct {"):
				currentType = "struct"
				d.Package.Structs = append(d.Package.Structs, Struct{
					Name:       strings.Split(line, " ")[1],
					Definition: line + "\n",
				})
				continue
			case strings.Contains(line, "interface {"):
				currentType = "interface"
				d.Package.Interfaces = append(d.Package.Interfaces, Interface{
					Name:       strings.Split(line, " ")[1],
					Definition: line + "\n",
				})
				continue
			default:
				currentType = "type"
				d.Package.Types = append(d.Package.Types, Type{
					Name:       strings.Split(line, " ")[1],
					Definition: line,
				})
			}
		case strings.HasPrefix(line, "var"):
			currentType = "var"
			d.Package.Variables = append(d.Package.Variables, Variable{
				Definition: line,
			})
			continue
		case strings.HasPrefix(line, "func ("):
			currentType = "func"
			if d.Package.getLastType() == nil {
				d.Package.Types = append(d.Package.Types, Type{})
			}
			*lastDefinition = line
			continue
		case strings.HasPrefix(line, "    "):
			currentType = "docs"
			if lastDocumentable != nil {
				lastDocumentable.addToDocs(line)
			}
			continue
		}
	}

	return nil
}

func parseVariable(input string) (v Variable) {
	if input == "" {
		return
	}
	input = strings.TrimSpace(input)
	v.Definition = input
	input = strings.TrimPrefix(input, "var ")
	input = strings.TrimPrefix(input, "const ")
	v.Name = strings.Split(input, " ")[0]
	input = strings.TrimSpace(strings.Join(strings.Split(input, " ")[1:], " "))
	if input == "" {
		return
	}
	if strings.Split(input, "")[0] == "=" {
		strings.TrimSpace(strings.Join(strings.Split(input, " ")[0:], " "))
	} else {
		v.Type = strings.TrimSpace(input)
	}
	return
}

func IsUpper(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsUpper(r) || !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func getFunctionName(input string) string {
	if !strings.HasPrefix(input, "func ") {
		return ""
	}

	if strings.HasPrefix(input, "func (") {
		return strings.Split(strings.TrimSpace(strings.Split(input, ")")[1]), "(")[0]
	}

	input = strings.TrimPrefix(input, "func ")
	return strings.Split(input, "(")[0]
}
