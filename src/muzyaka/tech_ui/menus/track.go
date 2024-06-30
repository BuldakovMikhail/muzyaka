package menus

import (
	"bufio"
	"fmt"
	"github.com/dixonwille/wmenu/v5"
	"github.com/pkg/errors"
	"log"
	"os"
	"src/internal/models/dto"
	"src/tech_ui/lib"
	"src/tech_ui/utils"
	"strings"
)

type trackWithClient struct {
	ClientEntity
	dto.TrackMeta
}

// TODO: скрыть бы это от внешнего мира

func (m *Menu) DownloadTrack(opt wmenu.Opt) error {
	item, ok := opt.Value.(trackWithClient)
	if !ok {
		log.Fatal("Could not cast option's value to Merch")
	}
	inputReader := bufio.NewReader(os.Stdin)

	fmt.Printf("Enter path to media: \n")
	path, _ := inputReader.ReadString('\n')
	path = strings.TrimRight(path, "\r\n")
	if path != "" {
		trackObject, err := utils.GetTrack(item.Client, item.Id, m.jwt)
		if err != nil {
			return err
		}

		err = lib.SaveFile(path, trackObject.Payload)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Menu) TrackActions(opt wmenu.Opt) error {
	item, ok := opt.Value.(trackWithClient)
	if !ok {
		log.Fatal("Could not cast option's value to Merch")
	}

	submenu := wmenu.NewMenu("Select option")
	submenu.Option("Download media", item, false, m.DownloadTrack)
	submenu.Option("Exit", nil, true, func(opt wmenu.Opt) error {
		return errExit
	})

	for {
		err := submenu.Run()
		fmt.Println()
		if err != nil {
			if errors.Is(err, errExit) {
				break
			}

			fmt.Printf("ERROR: %v\n\n", err)
		}
	}

	return nil
}

func (m *Menu) FindTracks(opt wmenu.Opt) error {
	var name string

	client, ok := opt.Value.(ClientEntity)

	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}

	inputReader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter track name:")
	name, _ = inputReader.ReadString('\n')
	name = strings.TrimRight(name, "\r\n")

	var curPage int = 1

	for {
		tracks, err := utils.FindTracks(client.Client, name, curPage, -1, m.jwt)
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
				fmt.Sprintf("Name: %s, Genre: %s, Source: %s", v.Name, genre, v.Source),
				trackWithClient{
					ClientEntity: client,
					TrackMeta:    *v,
				},
				false,
				m.TrackActions,
			)
		}
		tracksSubmenu.Option("Next", nil, false, func(opt wmenu.Opt) error {
			curPage++
			return nil
		})
		tracksSubmenu.Option("Exit", nil, true, func(opt wmenu.Opt) error {
			return errExit
		})

		err = tracksSubmenu.Run()
		if errors.Is(err, errExit) {
			break
		} else if err != nil {
			fmt.Printf("ERROR: %v\n\n", err)
		}
	}

	return nil
}
