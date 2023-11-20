package main

import (
	"2023_2_Holi/auth"
)

//	@title			Netfilx auth API
//	@version		1.0
//	@description	API of the nelfix auth service

//	@contact.name	Alex Chinaev
//	@contact.url	https://vk.com/l.chinaev
//	@contact.email	ax.chinaev@yandex.ru

//	@license.name	AS IS (NO WARRANTY)

// @host		127.0.0.1
// @schemes	http
// @BasePath	/api/v1/auth/
func main() {
	auth.StartService()
}
