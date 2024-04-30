package Global

import "fmt"

type UserInfo struct {
	ID     string
	Nombre string
}

func PrintUser(usr UserInfo) {
	fmt.Println("ID: " + usr.ID)
	fmt.Println("Nombre: " + usr.Nombre)
}

var Usuario UserInfo
