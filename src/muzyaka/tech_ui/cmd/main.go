package main

import (
	"net/http"
	"src/tech_ui/menus"
)

// TODO: нельзя создать плейлист без обложки
// TODO: операции с треками в плейлистах отличаются от тех что в лайках
// TODO: мб внести удаление треков из плейлиста в трекс экшен
// TODO: мб выводить ошибки красивее чем эту страшилку от бд
// TODO: если ввести два пробела между именами файлов, то ошибка буде неверный формат, мб решить как-то
// TODO: можно заменить названия пробелами, это плохо
// TODO: при добавлении альбома без треков выводит ошибку по которой ничего не понятно
// TODO: update album не получается обойтись без фото
// TODO: вообще почистить все от неиспользуемого кода
// TODO: удаление всех треков из альбома, альбом остается во вкладке удаления, до выхода из нее мб не страшно
func main() {
	menu := menus.NewMenu()
	menu.RunAuthMenu(http.DefaultClient)
}
