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
	TextInput   textinput.Model
	InputMode   bool
	TableResult table.Model
	TableMode   bool
}
