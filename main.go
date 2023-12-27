package main

import (
	api "docstruction/getconstructionmaterial/API"
	_ "docstruction/getconstructionmaterial/GCalls"
	server "docstruction/getconstructionmaterial/Server"
	"fmt"
)

func main() {
	err := api.RefreshPushNotificationWatch()
	if err != nil {
		fmt.Printf("Push Notification Refresh Error: %v", err)
	}
	server.Idle()
}
