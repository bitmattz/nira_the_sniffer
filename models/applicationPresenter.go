package models

type ApplicationPresenter struct {
	Choices  []string
	Cursor   int
	Selected map[int]struct{}
	Page     int
}
