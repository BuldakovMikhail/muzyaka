package menus

import (
	"github.com/dixonwille/wmenu/v5"
	"github.com/pkg/errors"
	"net/http"
)

type ClientEntity struct {
	Client *http.Client
}

var (
	errExit = errors.New("exiting")
)

type Menu struct {
	mainMenu     *wmenu.Menu
	musicianMenu *wmenu.Menu
	userMenu     *wmenu.Menu
	jwt          string
	id           uint64
	role         string
	musicianId   uint64
}

func NewMenu() *Menu {
	return &Menu{}
}
