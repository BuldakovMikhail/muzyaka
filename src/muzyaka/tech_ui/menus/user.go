package menus

import (
	"errors"
	"fmt"
	"github.com/dixonwille/wmenu/v5"
	"log"
	"src/tech_ui/utils"
)

func (m *Menu) GetMyLikedTracks(opt wmenu.Opt) error {
	client, ok := opt.Value.(ClientEntity)

	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}

	for {
		tracks, err := utils.GetLikedTracks(client.Client, m.id, m.jwt)
		if err != nil {
			return err
		}

		tracksSubmenu := wmenu.NewMenu("Select track")
		for _, v := range tracks {
			genre := "None"
			if v.Genre != nil {
				genre = *v.Genre
			}

			tracksSubmenu.Option(
				fmt.Sprintf("Name: %s, Genre: %s", v.Name, genre),
				TrackWithClient{
					ClientEntity: client,
					TrackMeta:    *v,
				},
				false,
				m.TrackActions,
			)
		}
		tracksSubmenu.Option("Exit", nil, true, func(opt wmenu.Opt) error {
			return errExit
		})

		err = tracksSubmenu.Run()
		if errors.Is(err, errExit) {
			break
		} else if err != nil {
			fmt.Printf("ERROR: %s\n\n", err.Error())
		}
	}

	return nil
}
