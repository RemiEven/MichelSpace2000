package assets

// Credit holds information about where an asset was found, who created it and how it is licensed
type Credit struct {
	Authors []string `json:"authors"`
	Source  string   `json:"source"`
	License string   `json:"license"`
}
