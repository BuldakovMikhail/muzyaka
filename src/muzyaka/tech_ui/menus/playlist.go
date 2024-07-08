package menus

import (
	"bufio"
	"fmt"
	"github.com/dixonwille/wmenu/v5"
	"github.com/pkg/errors"
	"log"
	"os"
	"src/internal/models"
	"src/internal/models/dto"
	"src/tech_ui/lib"
	"src/tech_ui/utils"
	"strings"
)

func (m *Menu) CreatePlaylist(opt wmenu.Opt) error {
	var name string
	var path string
	var description string

	client, ok := opt.Value.(ClientEntity)

	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}

	inputReader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter name:")
	name, _ = inputReader.ReadString('\n')
	name = strings.TrimRight(name, "\r\n")
	name = strings.Trim(name, " ")

	fmt.Println("Enter description:")
	description, _ = inputReader.ReadString('\n')
	description = strings.TrimRight(description, "\r\n")

	fmt.Println("Enter path to cover (*.png):")
	path, _ = inputReader.ReadString('\n')
	path = strings.TrimRight(path, "\r\n")
	if !lib.IsPNGFormat(path) {
		return models.ErrInvalidFileFormat
	}

	var arrOfBytes []byte
	if path != "" {
		var err error
		arrOfBytes, err = lib.ReadFile(path)
		if err != nil {
			return err
		}
	}

	err := utils.CreatePlaylist(
		client.Client,
		dto.PlaylistWithoutId{
			Name:        name,
			CoverFile:   arrOfBytes,
			Description: description,
		},
		m.id,
		m.jwt)

	return err
}

func (m *Menu) DeleteTrackFromPlaylist(opt wmenu.Opt) error {
	client, ok := opt.Value.(ClientEntity)

	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}

	items, err := utils.GetAllPlaylists(client.Client, m.id, m.jwt)
	if err != nil {
		return err
	}

	submenu := wmenu.NewMenu("Select playlist: ")
	for _, v := range items {
		submenu.Option(fmt.Sprintf("Name: %s, Description: %s", v.Name, v.Description),
			*v,
			false,
			func(opt wmenu.Opt) error {
				playlist, ok := opt.Value.(dto.Playlist)
				if !ok {
					log.Fatal("Could not cast option's value to Album")
				}

				tracks, err := utils.GetAllTracksFromPlaylist(client.Client, playlist.Id, m.jwt)
				if err != nil {
					return err
				}

				tracksSubmenu := wmenu.NewMenu("Select track for delete: ")
				for _, t := range tracks {
					genre := "None"
					if t.Genre != nil {
						genre = *t.Genre
					}
					tracksSubmenu.Option(
						fmt.Sprintf("Name: %s, Genre: %s", t.Name, genre),
						*t,
						false,
						func(opt wmenu.Opt) error {
							item, ok := opt.Value.(dto.TrackMeta)
							if !ok {
								log.Fatal("Could not cast option's value to Album")
							}

							err := utils.DeleteTrackFromPlaylist(client.Client, playlist.Id, item.Id, m.jwt)
							return err
						})
				}
				tracksSubmenu.Option("Exit", nil, true, func(opt wmenu.Opt) error {
					return nil
				})
				tracksSubmenu.Run()
				return nil
			})
	}
	submenu.Option("Exit", nil, true, func(_ wmenu.Opt) error {
		return errExit
	})

	for {
		err := submenu.Run()
		fmt.Println()
		if err != nil {
			if errors.Is(err, errExit) {
				break
			}

			fmt.Printf("ERROR: %s\n\n", err.Error())
		}
	}

	return nil
}

func (m *Menu) GetAllMyPlaylists(opt wmenu.Opt) error {
	client, ok := opt.Value.(ClientEntity)

	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}

	items, err := utils.GetAllPlaylists(client.Client, m.id, m.jwt)
	if err != nil {
		return err
	}

	submenu := wmenu.NewMenu("Open item: ")
	for _, v := range items {
		submenu.Option(fmt.Sprintf("Name: %s, Description: %s", v.Name, v.Description),
			*v,
			false,
			func(opt wmenu.Opt) error {
				item, ok := opt.Value.(dto.Playlist)
				if !ok {
					log.Fatal("Could not cast option's value to Merch")
				}

				inputReader := bufio.NewReader(os.Stdin)

				fmt.Printf("ID: %d\n", item.Id)
				fmt.Printf("Name: %s\n", item.Name)
				fmt.Printf("Description: %s\n", item.Description)

				fmt.Printf("Enter path for saving photo: \n")
				path, _ := inputReader.ReadString('\n')
				path = strings.TrimRight(path, "\r\n")
				if path != "" {
					err := lib.SaveFile(path, item.CoverFile)
					if err != nil {
						return err
					}
				}

				tracks, err := utils.GetAllTracksFromPlaylist(client.Client, item.Id, m.jwt)
				if err != nil {
					return err
				}

				submenuTracks := wmenu.NewMenu("Open track: ")

				for _, t := range tracks {
					genre := "None"
					if t.Genre != nil {
						genre = *t.Genre
					}

					submenuTracks.Option(
						fmt.Sprintf("Name: %s, Genre: %s", t.Name, genre),
						TrackWithClient{
							ClientEntity: client,
							TrackMeta:    *t,
						},
						false,
						m.TrackActions,
					)
				}
				submenuTracks.Option("Exit", nil, true, func(_ wmenu.Opt) error {
					return errExit
				})

				for {
					err := submenuTracks.Run()
					fmt.Println()
					if err != nil {
						if errors.Is(err, errExit) {
							break
						}

						fmt.Printf("ERROR: %s\n\n", err.Error())
					}
				}
				return nil
			})
	}
	submenu.Option("Exit", nil, true, func(_ wmenu.Opt) error {
		return errExit
	})

	for {
		err := submenu.Run()
		fmt.Println()
		if err != nil {
			if errors.Is(err, errExit) {
				break
			}

			fmt.Printf("ERROR: %s\n\n", err.Error())
		}
	}

	return nil
}

func (m *Menu) UpdatePlaylist(opt wmenu.Opt) error {
	client, ok := opt.Value.(ClientEntity)

	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}

	items, err := utils.GetAllPlaylists(client.Client, m.id, m.jwt)
	if err != nil {
		return err
	}

	submenu := wmenu.NewMenu("Update album: ")
	for _, v := range items {
		submenu.Option(fmt.Sprintf("Name: %s, Description: %s", v.Name, v.Description),
			*v,
			false,
			func(opt wmenu.Opt) error {
				var name string
				var path string
				var description string

				item, ok := opt.Value.(dto.Playlist)
				if !ok {
					log.Fatal("Could not cast option's value to Merch")
				}

				inputReader := bufio.NewReader(os.Stdin)

				fmt.Println("Enter name:")
				name, _ = inputReader.ReadString('\n')
				name = strings.TrimRight(name, "\r\n")
				name = strings.Trim(name, " ")

				fmt.Println("Enter description:")
				description, _ = inputReader.ReadString('\n')
				description = strings.TrimRight(description, "\r\n")

				fmt.Println("Enter path to cover (*.png):")
				path, _ = inputReader.ReadString('\n')
				path = strings.TrimRight(path, "\r\n")
				if !lib.IsPNGFormat(path) {
					return models.ErrInvalidFileFormat
				}

				var arrOfBytes []byte
				if path != "" {
					var err error
					arrOfBytes, err = lib.ReadFile(path)
					if err != nil {
						return err
					}
				}

				err = utils.UpdatePlaylist(client.Client,
					dto.PlaylistWithoutId{
						Name:        name,
						CoverFile:   arrOfBytes,
						Description: description,
					},
					item.Id,
					m.jwt,
				)

				return nil
			})
	}
	submenu.Option("Exit", nil, true, func(_ wmenu.Opt) error {
		return errExit
	})

	err = submenu.Run()
	if err != nil {
		return err
	}

	return nil
}

func (m *Menu) DeletePlaylist(opt wmenu.Opt) error {
	client, ok := opt.Value.(ClientEntity)

	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}

	items, err := utils.GetAllPlaylists(client.Client, m.id, m.jwt)
	if err != nil {
		return err
	}

	submenu := wmenu.NewMenu("Delete playlist: ")
	for _, v := range items {
		submenu.Option(fmt.Sprintf("Name: %s, Description: %s", v.Name, v.Description),
			*v,
			false,
			func(opt wmenu.Opt) error {
				item, ok := opt.Value.(dto.Playlist)
				if !ok {
					log.Fatal("Could not cast option's value to Album")
				}
				err = utils.DeletePlaylist(client.Client, item.Id, m.jwt)
				return nil
			})
	}
	submenu.Option("Exit", nil, true, func(_ wmenu.Opt) error {
		return nil
	})

	err = submenu.Run()
	if err != nil {
		return err
	}

	return nil
}
