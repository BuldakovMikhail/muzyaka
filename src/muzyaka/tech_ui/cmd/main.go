package main

import (
	"net/http"
	"src/tech_ui/menus"
)

func main() {
	menu := menus.NewMenu()
	menu.RunAuthMenu(http.DefaultClient)
}
