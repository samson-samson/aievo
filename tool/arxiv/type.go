package arxiv

// Result defines a search query result type.
type Result struct {
	Title       string
	Authors     []string
	Summary     string
	PdfURL      string
	PublishedAt string
}
