package menus

import (
	"fmt"
	"github.com/dixonwille/wmenu/v5"
	"github.com/pkg/errors"
	"net/http"
)

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
