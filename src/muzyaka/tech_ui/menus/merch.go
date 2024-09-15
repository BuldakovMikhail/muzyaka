package menus

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/dixonwille/wmenu/v5"
	"log"
	"os"
	"src/internal/models"
	"src/internal/models/dto"
	"src/tech_ui/lib"
	"src/tech_ui/utils"
	"strings"
)

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
	name = strings.Trim(name, " ")

	fmt.Println("Enter order URL:")
	orderUrl, _ = inputReader.ReadString('\n')
	orderUrl = strings.TrimRight(orderUrl, "\r\n")
	orderUrl = strings.Trim(orderUrl, "\r\n")

	fmt.Println("Enter paths to photos, separated by space (*.png):")
	paths, _ = inputReader.ReadString('\n')
	paths = strings.TrimRight(paths, "\r\n")

	var arrOfBytes [][]byte
	if paths != "" {
		var err error
		arrOfPaths := strings.Split(paths, " ")
		for _, v := range arrOfPaths {
			if !lib.IsPNGFormat(v) {
				return models.ErrInvalidFileFormat
			}
		}
		arrOfBytes, err = lib.ReadAllFilesFromArray(arrOfPaths)
		if err != nil {
			return err
		}
	}

	fmt.Println("Enter description:")
	description, _ = inputReader.ReadString('\n')
	description = strings.TrimRight(description, "\r\n")

	err := utils.CreateMerch(
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

func (m *Menu) GetMyMerch(opt wmenu.Opt) error {
	client, ok := opt.Value.(ClientEntity)

	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}

	items, err := utils.GetAllMerch(client.Client, m.musicianId, m.jwt)
	if err != nil {
		return err
	}

	submenu := wmenu.NewMenu("Open item: ")
	for _, v := range items {
		submenu.Option(fmt.Sprintf("Name: %s, URL: %s", v.Name, v.OrderUrl),
			*v,
			false,
			func(opt wmenu.Opt) error {
				item, ok := opt.Value.(dto.Merch)
				if !ok {
					log.Fatal("Could not cast option's value to Merch")
				}

				inputReader := bufio.NewReader(os.Stdin)

				fmt.Printf("ID: %d\n", item.Id)
				fmt.Printf("Name: %s\n", item.Name)
				fmt.Printf("URL: %s\n", item.OrderUrl)
				fmt.Printf("Description: %s\n", item.Description)

				for i, v := range item.PhotoFiles {
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

func (m *Menu) UpdateMerch(opt wmenu.Opt) error {
	client, ok := opt.Value.(ClientEntity)

	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}

	items, err := utils.GetAllMerch(client.Client, m.musicianId, m.jwt)
	if err != nil {
		return err
	}

	submenu := wmenu.NewMenu("Update item: ")
	for _, v := range items {
		submenu.Option(fmt.Sprintf("Name: %s, URL: %s", v.Name, v.OrderUrl),
			*v,
			false,
			func(opt wmenu.Opt) error {
				var name string
				var paths string
				var description string
				var orderUrl string

				item, ok := opt.Value.(dto.Merch)
				if !ok {
					log.Fatal("Could not cast option's value to Merch")
				}

				inputReader := bufio.NewReader(os.Stdin)

				fmt.Println("Enter name:")
				name, _ = inputReader.ReadString('\n')
				name = strings.TrimRight(name, "\r\n")
				name = strings.Trim(name, " ")

				fmt.Println("Enter order URL:")
				orderUrl, _ = inputReader.ReadString('\n')
				orderUrl = strings.TrimRight(orderUrl, "\r\n")
				orderUrl = strings.Trim(orderUrl, "\r\n")

				fmt.Println("Enter paths to photos, separated by space (*.png):")
				paths, _ = inputReader.ReadString('\n')
				paths = strings.TrimRight(paths, "\r\n")

				var arrOfBytes [][]byte
				if paths != "" {
					var err error
					arrOfPaths := strings.Split(paths, " ")
					for _, v := range arrOfPaths {
						if !lib.IsPNGFormat(v) {
							return models.ErrInvalidFileFormat
						}
					}
					arrOfBytes, err = lib.ReadAllFilesFromArray(arrOfPaths)
					if err != nil {
						return err
					}
				}

				fmt.Println("Enter description:")
				description, _ = inputReader.ReadString('\n')
				description = strings.TrimRight(description, "\r\n")

				err = utils.UpdateMerch(client.Client, dto.MerchWithoutId{
					Name:        name,
					PhotoFiles:  arrOfBytes,
					Description: description,
					OrderUrl:    orderUrl,
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

func (m *Menu) DeleteMerch(opt wmenu.Opt) error {
	client, ok := opt.Value.(ClientEntity)

	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}

	items, err := utils.GetAllMerch(client.Client, m.musicianId, m.jwt)
	if err != nil {
		return err
	}

	submenu := wmenu.NewMenu("Update item: ")
	for _, v := range items {
		submenu.Option(fmt.Sprintf("Name: %s, URL: %s", v.Name, v.OrderUrl),
			*v,
			false,
			func(opt wmenu.Opt) error {
				item, ok := opt.Value.(dto.Merch)
				if !ok {
					log.Fatal("Could not cast option's value to Merch")
				}
				err = utils.DeleteMerch(client.Client, item.Id, m.jwt)
				return err
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

func (m *Menu) FindMerch(opt wmenu.Opt) error {
	var name string

	client, ok := opt.Value.(ClientEntity)

	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}

	inputReader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter merch name:")
	name, _ = inputReader.ReadString('\n')
	name = strings.TrimRight(name, "\r\n")
	name = strings.Trim(name, " ")

	var curPage int = 1

	for {
		merch, err := utils.FindMerch(client.Client, name, curPage, -1, m.jwt)
		if err != nil {
			return err
		}

		submenu := wmenu.NewMenu("Select merch")
		for _, v := range merch {
			submenu.Option(fmt.Sprintf("Name: %s, URL: %s", v.Name, v.OrderUrl),
				*v,
				false,
				func(opt wmenu.Opt) error {
					item, ok := opt.Value.(dto.Merch)
					if !ok {
						log.Fatal("Could not cast option's value to Merch")
					}

					inputReader := bufio.NewReader(os.Stdin)

					fmt.Printf("ID: %d\n", item.Id)
					fmt.Printf("Name: %s\n", item.Name)
					fmt.Printf("URL: %s\n", item.OrderUrl)
					fmt.Printf("Description: %s\n", item.Description)

					for i, v := range item.PhotoFiles {
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
				})
		}
		if len(merch) != 0 {
			submenu.Option("Next", nil, false, func(opt wmenu.Opt) error {
				curPage++
				return nil
			})
		}
		if curPage > 1 {
			submenu.Option("Prev", nil, false, func(opt wmenu.Opt) error {
				curPage--
				return nil
			})
		}

		submenu.Option("Exit", nil, true, func(opt wmenu.Opt) error {
			return errExit
		})

		err = submenu.Run()
		if errors.Is(err, errExit) {
			break
		} else if err != nil {
			fmt.Printf("ERROR: %s\n\n", err.Error())
		}
	}

	return nil
}
