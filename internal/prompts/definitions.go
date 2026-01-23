package prompts

import (
	"text/template"
)

// PromptDefinition definition of an MCP prompt
type PromptDefinition struct {
	Name        string
	Description string
	Arguments   []PromptArgument
	FilePath    string
	Template    *template.Template
}

// PromptArgument definition of an MCP prompt argument
type PromptArgument struct {
	Name        string
	Description string
	Required    bool
}
