package main

import (
	api "docstruction/getconstructionmaterial/API"
	_ "docstruction/getconstructionmaterial/GCalls"
	server "docstruction/getconstructionmaterial/Server"
)

func main() {
	api.RefreshPushNotificationWatch()
	server.Idle()
}
