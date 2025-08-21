package types

// CategoryTool represents a tool within a category
type CategoryTool struct {
	Name        string
	Description string
	URL         string
	Icon        string
	Tags        []string
	IsNew       bool
	IsPopular   bool
}

// CategoryData represents the data for a category page
type CategoryData struct {
	Name        string
	Description string
	Icon        string
	ToolCount   string
	Tools       []CategoryTool
	SearchHint  string
	Color       string
}
