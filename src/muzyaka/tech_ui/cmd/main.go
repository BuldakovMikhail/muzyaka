package main

import (
	"net/http"
	"src/tech_ui/menus"
)

// TODO: мб выводить ошибки красивее чем эту страшилку от бд
// TODO: при добавлении альбома без треков выводит ошибку по которой ничего не понятно
// TODO: вообще почистить все от неиспользуемого кода

func main() {
	menu := menus.NewMenu()
	menu.RunAuthMenu(http.DefaultClient)
}
