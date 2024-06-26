package menus

import (
	"bufio"
	"fmt"
	"github.com/dixonwille/wmenu/v5"
	"log"
	"os"
	"src/internal/models/dto"
	"src/tech_ui/lib"
	"src/tech_ui/utils"
	"strings"
)

func (m *Menu) GetMyMusicianProfile(opt wmenu.Opt) error {
	client, ok := opt.Value.(ClientEntity)

	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}

	profile, err := utils.GetMusician(client.Client, m.musicianId, m.jwt)
	if err != nil {
		return err
	}

	inputReader := bufio.NewReader(os.Stdin)

	fmt.Printf("ID: %d\n", profile.Id)
	fmt.Printf("Name: %s\n", profile.Name)
	fmt.Printf("Description: %s\n", profile.Description)

	for i, v := range profile.PhotoFiles {
		fmt.Printf("Enter path to photo №%d: \n", i+1)
		path, _ := inputReader.ReadString('\n')
		path = strings.TrimRight(path, "\r\n")
		if path != "" {
			err := lib.SaveFile(path, v)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// TODO: при обновлении музыканта нужно дать возможность пропускать изменение фотографии

func (m *Menu) UpdateMyMusicianProfile(opt wmenu.Opt) error {
	var name string
	var paths string
	var description string

	client, ok := opt.Value.(ClientEntity)

	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}

	inputReader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter name:")
	name, _ = inputReader.ReadString('\n')
	name = strings.TrimRight(name, "\r\n")

	fmt.Println("Enter paths to photos, separated by space:")
	paths, _ = inputReader.ReadString('\n')

	arrOfPaths := strings.Split(paths, " ")
	arrOfBytes, err := lib.ReadAllFilesFromArray(arrOfPaths)
	if err != nil {
		return err
	}

	fmt.Println("Enter description:")
	description, _ = inputReader.ReadString('\n')
	description = strings.TrimRight(description, "\r\n")

	err = utils.UpdateMusician(client.Client, dto.MusicianWithoutId{
		Name:        name,
		PhotoFiles:  arrOfBytes,
		Description: description,
	}, m.musicianId, m.jwt)

	return nil
}
