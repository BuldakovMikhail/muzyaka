package menus

import (
	"fmt"
	"github.com/dixonwille/wmenu/v5"
	"github.com/pkg/errors"
	"net/http"
)

func (m *Menu) AddOptionsUser(client *http.Client) {
	m.userMenu.Option("Find Track", ClientEntity{client}, false, m.FindTracks)
	m.userMenu.Option("Get my favorites", ClientEntity{client}, false, m.GetMyLikedTracks)
	m.userMenu.Option("Exit", ClientEntity{client}, false, func(_ wmenu.Opt) error {
		return errExit
	})
}

func (m *Menu) RunUserMenu(client *http.Client) error {
	m.userMenu = wmenu.NewMenu("Enter your option:")
	m.AddOptionsUser(client)

	for {
		err := m.userMenu.Run()
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
