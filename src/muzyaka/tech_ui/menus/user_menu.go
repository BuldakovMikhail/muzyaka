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
	m.userMenu.Option("Create playlist", ClientEntity{client}, false, m.CreatePlaylist)
	m.userMenu.Option("Get all my playlists", ClientEntity{client}, false, m.GetAllMyPlaylists)
	m.userMenu.Option("Update my playlist", ClientEntity{client}, false, m.UpdatePlaylist)
	m.userMenu.Option("Delete my playlist", ClientEntity{client}, false, m.DeletePlaylist)
	m.userMenu.Option("Delete tracks from playlist", ClientEntity{client}, false, m.DeleteTrackFromPlaylist)
	m.userMenu.Option("Find Merch", ClientEntity{client}, false, m.FindMerch)
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

			fmt.Printf("ERROR: %s\n\n", err.Error())
		}
	}

	fmt.Printf("Exited menu.\n")
	return nil
}
