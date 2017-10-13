package teamwork

// Currently supports: Int and String
type Pages struct {
	Page    int `header:"X-Page"`
	Pages   int `header:"X-Pages"`
	Records int `header:"X-Records"`
}
