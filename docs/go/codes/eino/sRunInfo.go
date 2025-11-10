type RunInfo struct {
	// Name is the graph node name for display purposes, not unique.
	// Passed from compose.WithNodeName().
	Name      string
	Type      string
	Component components.Component
}