package menus

import (
	"github.com/dixonwille/wmenu/v5"
	"net/http"
)

type ClientEntity struct {
	Client *http.Client
}

type Menu struct {
	mainMenu   *wmenu.Menu
	jwt        string
	id         uint64
	role       string
	musicianId uint64
}

func NewMenu() *Menu {
	return &Menu{}
}
