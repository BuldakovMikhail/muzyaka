package menus

import (
	"fmt"
	"github.com/dixonwille/wmenu/v5"
	"github.com/pkg/errors"
	"net/http"
)

// TODO: мб нет фото?
func (m *Menu) AddOptionsMusician(client *http.Client) {
	m.musicianMenu.Option("Add Merch", ClientEntity{client}, false, m.CreateMerch)
	m.musicianMenu.Option("Get all my merch", ClientEntity{client}, false, m.GetMyMerch)
	m.musicianMenu.Option("Update my merch", ClientEntity{client}, false, m.UpdateMerch)
	m.musicianMenu.Option("Exit", ClientEntity{client}, false, func(_ wmenu.Opt) error {
		return errExit
	})
}

func (m *Menu) RunMusicianMenu(client *http.Client) error {
	m.musicianMenu = wmenu.NewMenu("Enter your option:")
	m.AddOptionsMusician(client)

	for {
		err := m.musicianMenu.Run()
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
