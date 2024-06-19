package menus

import (
	"bufio"
	"fmt"
	"github.com/dixonwille/wmenu/v5"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"os"
	"src/internal/models/dto"
	"src/tech_ui/lib"
	"src/tech_ui/utils"
	"strings"
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

func (m *Menu) CreateMerch(opt wmenu.Opt) error {
	var name string
	var paths string
	var description string
	var orderUrl string

	client, ok := opt.Value.(ClientEntity)

	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}

	inputReader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter name:")
	name, _ = inputReader.ReadString('\n')
	name = strings.TrimRight(name, "\r\n")

	fmt.Println("Enter order URL:")
	orderUrl, _ = inputReader.ReadString('\n')
	orderUrl = strings.TrimRight(orderUrl, "\r\n")

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

	err = utils.CreateMerch(
		client.Client,
		dto.MerchWithoutId{
			Name:        name,
			PhotoFiles:  arrOfBytes,
			Description: description,
			OrderUrl:    orderUrl,
		},
		m.musicianId,
		m.jwt)

	if err != nil {
		return err
	}

	return nil
}
