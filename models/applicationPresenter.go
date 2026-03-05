package models

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
)

type ApplicationPresenter struct {
	Choices     []string
	Cursor      int
	Selected    map[int]struct{}
	Page        int
	InputMode   bool
	TextInput   textinput.Model
	TableMode   bool
	TableResult table.Model
	IsLoading   bool
}
