package data

type Page struct {
	Title string
	Body  []byte
}

type PageStore interface {
	SavePage(page *Page) error
	LoadPage(title string) (*Page, error)
}
