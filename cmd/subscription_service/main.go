package main

import (
	"2023_2_Holi/subscription"
)

//	@title			Netfilx API
//	@version		1.0
//	@description	API of the nelfix film and series service

//	@contact.name	Aleksej Moldovanov
//	@contact.url	https://vk.com/yepkekw
//	@contact.email	3592703@gmail.com

//	@license.name	AS IS (NO WARRANTY)

// @host			127.0.0.1
// @schemes			http
// @BasePath		/api/v1/films/ & /api/v1/series/
func main() {
	subscription.StartService()
}
