package menus

import (
	"fmt"
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
	mainMenu   *wmenu.Menu
	jwt        string
	id         uint64
	role       string
	musicianId uint64
}

func NewMenu() *Menu {
	return &Menu{}
}

func (m *Menu) AddOptionsMain(client *http.Client) {
	m.mainMenu.Option("Sign up as user", ClientEntity{client}, false, m.SignUpAsUser)
	m.mainMenu.Option("Sign up as musician", ClientEntity{client}, false, m.SignUpAsMusician)
	m.mainMenu.Option("Sign in", ClientEntity{client}, false, m.SignIn)
	m.mainMenu.Option("Exit", ClientEntity{client}, false, func(_ wmenu.Opt) error {
		return errExit
	})
}

func (m *Menu) RunMenu(client *http.Client) error {
	m.mainMenu = wmenu.NewMenu("Please select options")
	m.AddOptionsMain(client)

	for {
		err := m.mainMenu.Run()
		fmt.Println()
		if err != nil {
			if errors.Is(err, errExit) {
				break
			}
			fmt.Printf("ERROR: %v\n\n", err)
		}

	}

	fmt.Printf("Exited menu.\n")

	return nil
}
