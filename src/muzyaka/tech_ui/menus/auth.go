package menus

import (
	"bufio"
	"fmt"
	"github.com/dixonwille/wmenu/v5"
	"log"
	"os"
	"src/internal/lib/api/response"
	"src/internal/models/dto"
	"src/tech_ui/lib"
	"src/tech_ui/utils"
	"strings"
)

func (m *Menu) SignIn(opt wmenu.Opt) error {
	var login string
	var password string

	client, ok := opt.Value.(ClientEntity)

	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}

	fmt.Println("Enter login:")
	fmt.Scan(&login)

	fmt.Println("Enter password:")
	fmt.Scan(&password)

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

	//switch m.role {
	//case usecase.UserRole:
	//	m.RunUserMenu(client.Client)
	//case usecase.MusicianRole:
	//	m.RunMusicianMenu(client.Client)
	//}

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

	fmt.Println("Enter login:")
	fmt.Scan(&login)

	fmt.Println("Enter password:")
	fmt.Scan(&password)

	fmt.Println("Enter name:")
	fmt.Scan(&name)

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

	//switch m.role {
	//case usecase.UserRole:
	//	m.RunUserMenu(client.Client)
	//case usecase.MusicianRole:
	//	m.RunMusicianMenu(client.Client)
	//}

	fmt.Println(response.StatusOK)

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

	//switch m.role {
	//case usecase.UserRole:
	//	m.RunUserMenu(client.Client)
	//case usecase.MusicianRole:
	//	m.RunMusicianMenu(client.Client)
	//}

	fmt.Println(response.StatusOK)

	return nil
}
