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

func (m *Menu) CreateAlbum(opt wmenu.Opt) error {
	var name string
	var path string
	var albumType string

	client, ok := opt.Value.(ClientEntity)

	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}

	inputReader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter name:")
	name, _ = inputReader.ReadString('\n')
	name = strings.TrimRight(name, "\r\n")

	fmt.Println("Enter path to cover:")
	path, _ = inputReader.ReadString('\n')
	path = strings.TrimRight(path, "\r\n")

	var arrOfBytes [][]byte
	if path != "" {
		var err error
		arrOfBytes, err = lib.ReadAllFilesFromArray([]string{path})
		if err != nil {
			return err
		}
	}

	submenu := wmenu.NewMenu("Select album type: ")
	submenu.Option("LP", nil, true, func(opt wmenu.Opt) error {
		albumType = "LP"
		return nil
	})
	submenu.Option("Single", nil, false, func(opt wmenu.Opt) error {
		albumType = "single"
		return nil
	})
	submenu.Option("EP", nil, false, func(opt wmenu.Opt) error {
		albumType = "EP"
		return nil
	})
	submenu.Run()

	var tracks []*dto.TrackObjectWithoutId

	submenuTracks := wmenu.NewMenu("Add tracks: ")
	submenuTracks.Option("Add track",
		nil,
		true,
		func(opt wmenu.Opt) error {
			inputReader := bufio.NewReader(os.Stdin)

			fmt.Println("Enter name:")
			trackName, _ := inputReader.ReadString('\n')
			trackName = strings.TrimRight(trackName, "\r\n")

			fmt.Println("Enter genre:")
			genre, _ := inputReader.ReadString('\n')
			genre = strings.TrimRight(genre, "\r\n")

			fmt.Println("Enter path to payload:")
			source, _ := inputReader.ReadString('\n')
			source = strings.TrimRight(source, "\r\n")

			var payload [][]byte
			if path != "" {
				var err error
				payload, err = lib.ReadAllFilesFromArray([]string{source})
				if err != nil {
					return err
				}
			}

			genreRef := &genre
			if genre == "" {
				genreRef = nil
			}

			tracks = append(tracks, &dto.TrackObjectWithoutId{
				TrackMetaWithoutId: dto.TrackMetaWithoutId{
					Source: source, // TODO: заменить на генерацию
					Name:   trackName,
					Genre:  genreRef,
				},
				Payload:     payload[0],
				PayloadSize: int64(len(payload[0])), // TODO: убрать приведение типов
			}) // TODO: не добавляется в outbox

			return nil
		})
	submenuTracks.Option("Exit", nil, false, func(_ wmenu.Opt) error {
		return errExit
	})

	for {
		err := submenuTracks.Run()
		fmt.Println()
		if err != nil {
			if errors.Is(err, errExit) {
				break
			}

			fmt.Printf("ERROR: %v\n\n", err)
		}
	}

	err := utils.CreateAlbum(
		client.Client,
		dto.AlbumWithTracks{
			AlbumWithoutId: dto.AlbumWithoutId{
				Name:      name,
				CoverFile: arrOfBytes[0],
				Type:      albumType,
			},
			Tracks: tracks,
		},
		m.musicianId,
		m.jwt)

	if err != nil {
		return err
	}

	return nil
}
