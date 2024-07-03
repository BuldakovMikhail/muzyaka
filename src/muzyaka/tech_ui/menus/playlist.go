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

	fmt.Println("Enter description:")
	description, _ = inputReader.ReadString('\n')
	description = strings.TrimRight(description, "\r\n")

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

	err := utils.CreatePlaylist(
		client.Client,
		dto.PlaylistWithoutId{
			Name:        name,
			CoverFile:   arrOfBytes[0],
			Description: description,
		},
		m.id,
		m.jwt)

	return err
}
