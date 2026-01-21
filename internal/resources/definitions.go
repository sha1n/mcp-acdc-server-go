package resources

// ResourceDefinition definition of an MCP resource
type ResourceDefinition struct {
	URI         string
	Name        string
	Description string
	MIMEType    string
	FilePath    string
	Keywords    []string // Optional keywords for search boosting
}
