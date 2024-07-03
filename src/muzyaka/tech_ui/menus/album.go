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
			if source != "" {
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

func (m *Menu) GetAllMyAlbums(opt wmenu.Opt) error {
	client, ok := opt.Value.(ClientEntity)

	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}
	// TODO: убрать вложенность
	// TODO: мета информация о треке передается 2 раза, мб в треке отдавать только пейлоад

	items, err := utils.GetAllAlbums(client.Client, m.musicianId, m.jwt)
	if err != nil {
		return err
	}

	submenu := wmenu.NewMenu("Open item: ")
	for _, v := range items {
		submenu.Option(fmt.Sprintf("Name: %s, Type: %s", v.Name, v.Type),
			*v,
			false,
			func(opt wmenu.Opt) error {
				item, ok := opt.Value.(dto.Album)
				if !ok {
					log.Fatal("Could not cast option's value to Merch")
				}

				inputReader := bufio.NewReader(os.Stdin)

				fmt.Printf("ID: %d\n", item.Id)
				fmt.Printf("Name: %s\n", item.Name)
				fmt.Printf("Type: %s\n", item.Type)

				fmt.Printf("Enter path to photo: \n")
				path, _ := inputReader.ReadString('\n')
				path = strings.TrimRight(path, "\r\n")
				if path != "" {
					err := lib.SaveFile(path, item.CoverFile)
					if err != nil {
						return err
					}
				}

				tracks, err := utils.GetAllTracks(client.Client, item.Id, m.jwt)
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
						fmt.Sprintf("Name: %s, Genre: %s, Source: %s", t.Name, genre, t.Source),
						*t,
						false,
						func(opt wmenu.Opt) error {
							item, ok := opt.Value.(dto.TrackMeta)
							if !ok {
								log.Fatal("Could not cast option's value to Merch")
							}

							inputReader := bufio.NewReader(os.Stdin)

							fmt.Printf("ID: %d\n", item.Id)
							fmt.Printf("Name: %s\n", item.Name)
							fmt.Printf("Genre: %s\n", item.Genre)
							fmt.Printf("Source: %s\n", item.Source)

							fmt.Printf("Enter path to media: \n")
							path, _ := inputReader.ReadString('\n')
							path = strings.TrimRight(path, "\r\n")
							if path != "" {
								trackObject, err := utils.GetTrack(client.Client, item.Id, m.jwt)
								if err != nil {
									return err
								}

								err = lib.SaveFile(path, trackObject.Payload)
								if err != nil {
									return err
								}
							}
							return nil
						})
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

						fmt.Printf("ERROR: %v\n\n", err)
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

			fmt.Printf("ERROR: %v\n\n", err)
		}
	}

	return nil
}

func (m *Menu) UpdateAlbum(opt wmenu.Opt) error {
	client, ok := opt.Value.(ClientEntity)

	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}

	items, err := utils.GetAllAlbums(client.Client, m.musicianId, m.jwt)
	if err != nil {
		return err
	}

	submenu := wmenu.NewMenu("Update album: ")
	for _, v := range items {
		submenu.Option(fmt.Sprintf("Name: %s, Type: %s", v.Name, v.Type),
			*v,
			false,
			func(opt wmenu.Opt) error {
				var name string
				var path string
				var albumType string

				item, ok := opt.Value.(dto.Album)
				if !ok {
					log.Fatal("Could not cast option's value to Merch")
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

				err = utils.UpdateAlbum(client.Client, dto.AlbumWithoutId{
					Name:      name,
					CoverFile: arrOfBytes[0],
					Type:      albumType,
				}, item.Id, m.jwt)

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

func (m *Menu) DeleteAlbum(opt wmenu.Opt) error {
	client, ok := opt.Value.(ClientEntity)

	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}

	items, err := utils.GetAllAlbums(client.Client, m.musicianId, m.jwt)
	if err != nil {
		return err
	}

	submenu := wmenu.NewMenu("Delete album: ")
	for _, v := range items {
		submenu.Option(fmt.Sprintf("Name: %s, Type: %s", v.Name, v.Type),
			*v,
			false,
			func(opt wmenu.Opt) error {
				item, ok := opt.Value.(dto.Album)
				if !ok {
					log.Fatal("Could not cast option's value to Album")
				}
				err = utils.DeleteAlbum(client.Client, item.Id, m.jwt)
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

func (m *Menu) AddTrackToAlbum(opt wmenu.Opt) error {
	client, ok := opt.Value.(ClientEntity)

	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}

	items, err := utils.GetAllAlbums(client.Client, m.musicianId, m.jwt)
	if err != nil {
		return err
	}

	submenu := wmenu.NewMenu("Select album: ")
	for _, v := range items {
		submenu.Option(fmt.Sprintf("Name: %s, Type: %s", v.Name, v.Type),
			*v,
			false,
			func(opt wmenu.Opt) error {
				item, ok := opt.Value.(dto.Album)
				if !ok {
					log.Fatal("Could not cast option's value to Album")
				}
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
				if source != "" {
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

				track := dto.TrackObjectWithoutId{
					TrackMetaWithoutId: dto.TrackMetaWithoutId{
						Source: source, // TODO: заменить на генерацию
						Name:   trackName,
						Genre:  genreRef,
					},
					Payload:     payload[0],
					PayloadSize: int64(len(payload[0])), // TODO: убрать приведение типов
				}

				err = utils.AddTrack(client.Client, track, item.Id, m.jwt)
				if err != nil {
					return err
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

			fmt.Printf("ERROR: %v\n\n", err)
		}
	}

	return nil
}

func (m *Menu) DeleteTrackFromAlbum(opt wmenu.Opt) error {
	client, ok := opt.Value.(ClientEntity)

	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}

	items, err := utils.GetAllAlbums(client.Client, m.musicianId, m.jwt)
	if err != nil {
		return err
	}

	submenu := wmenu.NewMenu("Select album: ")
	for _, v := range items {
		submenu.Option(fmt.Sprintf("Name: %s, Type: %s", v.Name, v.Type),
			*v,
			false,
			func(opt wmenu.Opt) error {
				item, ok := opt.Value.(dto.Album)
				if !ok {
					log.Fatal("Could not cast option's value to Album")
				}

				tracks, err := utils.GetAllTracks(client.Client, item.Id, m.jwt)
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
						fmt.Sprintf("Name: %s, Source: %s, Genre: %s", t.Name, t.Source, genre),
						*t,
						false,
						func(opt wmenu.Opt) error {
							item, ok := opt.Value.(dto.TrackMeta)
							if !ok {
								log.Fatal("Could not cast option's value to Album")
							}

							err := utils.DeleteTrack(client.Client, item.Id, m.jwt)
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

			fmt.Printf("ERROR: %v\n\n", err)
		}
	}

	return nil
}
