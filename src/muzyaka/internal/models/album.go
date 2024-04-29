package models

// TODO: replace cover with array of bytes
// Вопрос заключается в том, как лучше разделить фото и сам альбом

type Album struct {
	Id    uint64
	Name  string
	Cover string
	Type  string
}
