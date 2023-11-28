package main

import "2023_2_Holi/profile"

//	@title			Netfilx auth API
//	@version		1.0
//	@description	API of the nelfix profile service

//	@contact.name	Alexander Krylov
//	@contact.url	https://vk.com/shelbyo0
//	@contact.email	krylov.sanches@mail.ru

//	@license.name	AS IS (NO WARRANTY)

// @host			127.0.0.1
// @schemes			http
// @BasePath		/api/v1/profile/
func main() {
	profile.StartService()
}
