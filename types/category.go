package types

// CategoryTool represents a tool within a tools
type CategoryTool struct {
	Name        string
	Description string
	URL         string
	Icon        string
	Tags        []string
	IsNew       bool
	IsPopular   bool
}

// CategoryData represents the data for a tools page
type CategoryData struct {
	Name        string
	Description string
	Icon        string
	ToolCount   string
	Tools       []CategoryTool
	SearchHint  string
	Color       string
}
