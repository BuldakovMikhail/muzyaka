package menus

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/dixonwille/wmenu/v5"
	"log"
	"net/http"
	"os"
	"src/internal/domain/auth/usecase"
	"src/internal/lib/api/response"
	"src/internal/models"
	"src/internal/models/dto"
	"src/tech_ui/lib"
	"src/tech_ui/utils"
	"strings"
)

func (m *Menu) AddOptionsMain(client *http.Client) {
	m.mainMenu.Option("Sign up as user", ClientEntity{client}, false, m.SignUpAsUser)
	m.mainMenu.Option("Sign up as musician", ClientEntity{client}, false, m.SignUpAsMusician)
	m.mainMenu.Option("Sign in", ClientEntity{client}, false, m.SignIn)
	m.mainMenu.Option("Exit", ClientEntity{client}, false, func(_ wmenu.Opt) error {
		return errExit
	})
}

func (m *Menu) RunAuthMenu(client *http.Client) error {
	m.mainMenu = wmenu.NewMenu("Please select options")
	m.AddOptionsMain(client)

	for {
		err := m.mainMenu.Run()
		fmt.Println()
		if err != nil {
			if errors.Is(err, errExit) {
				break
			}
			fmt.Printf("ERROR: %s\n", err.Error())
		}

	}

	fmt.Printf("Exited menu.\n")

	return nil
}

func (m *Menu) SignIn(opt wmenu.Opt) error {
	var login string
	var password string

	client, ok := opt.Value.(ClientEntity)

	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}

	inputReader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter login:")
	login, _ = inputReader.ReadString('\n')
	login = strings.TrimRight(login, "\r\n")

	fmt.Println("Enter password:")
	password, _ = inputReader.ReadString('\n')
	password = strings.TrimRight(password, "\r\n")

	jwt, err := utils.SignIn(client.Client, dto.SignIn{
		Email:    login,
		Password: password,
	})
	if err != nil {
		return err
	}

	m.jwt = jwt
	getMe, err := utils.GetMe(client.Client, m.jwt)
	if err != nil {
		return err
	}

	m.id = getMe.UserId
	m.role = getMe.Role
	m.musicianId = getMe.MusicianId

	switch m.role {
	case usecase.UserRole:
		m.RunUserMenu(client.Client)
	case usecase.MusicianRole:
		m.RunMusicianMenu(client.Client)
	}

	fmt.Println(response.StatusOK)

	return nil
}

func (m *Menu) SignUpAsUser(opt wmenu.Opt) error {
	var login string
	var password string
	var name string

	client, ok := opt.Value.(ClientEntity)

	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}

	inputReader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter login:")
	login, _ = inputReader.ReadString('\n')
	login = strings.TrimRight(login, "\r\n")

	fmt.Println("Enter password:")
	password, _ = inputReader.ReadString('\n')
	password = strings.TrimRight(password, "\r\n")

	fmt.Println("Enter name:")
	name, _ = inputReader.ReadString('\n')
	name = strings.TrimRight(name, "\r\n")
	name = strings.Trim(name, " ")

	jwt, err := utils.SignUpAsUser(client.Client,
		dto.SignUp{
			dto.UserInfo{
				Name:     name,
				Password: password,
				Email:    login,
			}})

	if err != nil {
		return err
	}

	m.jwt = jwt
	getMe, err := utils.GetMe(client.Client, m.jwt)
	if err != nil {
		return err
	}

	m.id = getMe.UserId
	m.role = getMe.Role
	m.musicianId = getMe.MusicianId

	fmt.Println(response.StatusOK)
	m.RunUserMenu(client.Client)

	return nil
}

func (m *Menu) SignUpAsMusician(opt wmenu.Opt) error {
	var login string
	var password string
	var name string
	var paths string
	var description string

	client, ok := opt.Value.(ClientEntity)

	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}

	inputReader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter login:")
	login, _ = inputReader.ReadString('\n')
	login = strings.TrimRight(login, "\r\n")

	fmt.Println("Enter password:")
	password, _ = inputReader.ReadString('\n')
	password = strings.TrimRight(password, "\r\n")

	fmt.Println("Enter name:")
	name, _ = inputReader.ReadString('\n')
	name = strings.TrimRight(name, "\r\n")
	name = strings.Trim(name, " ")

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

	jwt, err := utils.SignUpAsMusician(client.Client, dto.SignUpMusician{
		UserInfo: dto.UserInfo{
			Name:     name,
			Password: password,
			Email:    login,
		},
		MusicianWithoutId: dto.MusicianWithoutId{
			Name:        name,
			PhotoFiles:  arrOfBytes,
			Description: description,
		},
	})
	if err != nil {
		return err
	}

	m.jwt = jwt
	getMe, err := utils.GetMe(client.Client, m.jwt)
	if err != nil {
		return err
	}

	m.id = getMe.UserId
	m.role = getMe.Role
	m.musicianId = getMe.MusicianId

	fmt.Println(response.StatusOK)
	m.RunMusicianMenu(client.Client)

	return nil
}
