package menus

import (
	"fmt"
	"github.com/dixonwille/wmenu/v5"
	"log"
	"src/internal/lib/api/response"
	"src/tech_ui/utils"
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

	jwt, err := utils.SignIn(client.Client, login, password)
	if err != nil {
		return err
	}

	m.jwt = jwt
	getMe, err := utils.GetMe(client.Client)
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
