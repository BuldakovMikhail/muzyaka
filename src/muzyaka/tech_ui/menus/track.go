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

type TrackWithClient struct {
	ClientEntity
	dto.TrackMeta
}

// TODO: скрыть бы это от внешнего мира

func (m *Menu) GetSameTracks(opt wmenu.Opt) error {
	item, ok := opt.Value.(TrackWithClient)

	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}

	var curPage int = 1
	for {
		tracks, err := utils.GetSameTracks(item.Client, item.Id, curPage, -1, m.jwt)
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
					ClientEntity: item.ClientEntity,
					TrackMeta:    *v,
				},
				false,
				m.TrackActions,
			)
		}
		if len(tracks) != 0 {
			tracksSubmenu.Option("Next", nil, false, func(opt wmenu.Opt) error {
				curPage++
				return nil
			})
		}
		if curPage > 1 {
			tracksSubmenu.Option("Prev", nil, false, func(opt wmenu.Opt) error {
				curPage--
				return nil
			})
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

func (m *Menu) DownloadTrack(opt wmenu.Opt) error {
	item, ok := opt.Value.(TrackWithClient)
	if !ok {
		log.Fatal("Could not cast option's value to Merch")
	}
	inputReader := bufio.NewReader(os.Stdin)

	fmt.Printf("Enter path for saving media: \n")
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

func (m *Menu) LikeTrack(opt wmenu.Opt) error {
	item, ok := opt.Value.(TrackWithClient)
	if !ok {
		log.Fatal("Could not cast option's value to Merch")
	}
	err := utils.LikeTrack(item.Client, dto.Like{TrackId: item.Id}, m.id, m.jwt)
	return err
}

func (m *Menu) AddTrackToPlaylist(opt wmenu.Opt) error {
	item, ok := opt.Value.(TrackWithClient)
	if !ok {
		log.Fatal("Could not cast option's value to Merch")
	}

	items, err := utils.GetAllPlaylists(item.Client, m.id, m.jwt)
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
					log.Fatal("Could not cast option's value to Merch")
				}

				err := utils.AddTrackToPlaylist(
					item.Client,
					playlist.Id,
					dto.AddTrackPlaylistRequest{TrackId: item.Id},
					m.jwt,
				)

				return err
			})
	}
	submenu.Option("Exit", nil, true, func(_ wmenu.Opt) error {
		return nil
	})

	err = submenu.Run()
	return err
}

func (m *Menu) TrackActions(opt wmenu.Opt) error {
	item, ok := opt.Value.(TrackWithClient)
	if !ok {
		log.Fatal("Could not cast option's value to Merch")
	}

	for {
		isLiked, err := utils.IsTrackLiked(item.Client, m.id, item.Id, m.jwt)
		if err != nil {
			return err
		}
		submenu := wmenu.NewMenu("Select option")
		submenu.Option("Download media", item, false, m.DownloadTrack)

		if !isLiked {
			submenu.Option("Like track", item, false, m.LikeTrack)
		} else {
			submenu.Option("Dislike track", item, false, m.DislikeTrack)
		}
		submenu.Option("Get same tracks", item, false, m.GetSameTracks)
		submenu.Option("Add track to playlist", item, false, m.AddTrackToPlaylist)
		submenu.Option("Exit", nil, true, func(opt wmenu.Opt) error {
			return errExit
		})

		err = submenu.Run()
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
	name = strings.Trim(name, " ")

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
				fmt.Sprintf("Name: %s, Genre: %s", v.Name, genre),
				TrackWithClient{
					ClientEntity: client,
					TrackMeta:    *v,
				},
				false,
				m.TrackActions,
			)
		}
		if len(tracks) != 0 {
			tracksSubmenu.Option("Next", nil, false, func(opt wmenu.Opt) error {
				curPage++
				return nil
			})
		}
		if curPage > 1 {
			tracksSubmenu.Option("Prev", nil, false, func(opt wmenu.Opt) error {
				curPage--
				return nil
			})
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

func (m *Menu) DislikeTrack(opt wmenu.Opt) error {
	item, ok := opt.Value.(TrackWithClient)
	if !ok {
		log.Fatal("Could not cast option's value to Merch")
	}
	err := utils.DislikeTrack(item.Client, dto.Dislike{TrackId: item.Id}, m.id, m.jwt)
	return err
}
