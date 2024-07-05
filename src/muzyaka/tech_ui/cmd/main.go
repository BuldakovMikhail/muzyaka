package main

import (
	"net/http"
	"src/tech_ui/menus"
)

// TODO: добавить проверку на формат трека при загрузке
// TODO: если есть next, то добавить prev
func main() {
	menu := menus.NewMenu()
	menu.RunAuthMenu(http.DefaultClient)
}
