package article

type Article struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Slug    string `json:"slug"`
	Summary string `json:"summary"`
	Content string `json:"content"`
	Status  string `json:"status"`
}
