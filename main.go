package main

import (
	api "docstruction/getconstructionmaterial/API"
	_ "docstruction/getconstructionmaterial/GCalls"
	server "docstruction/getconstructionmaterial/Server"
	"fmt"
	"time"
)

func main() {
	err := api.RefreshPushNotificationWatch()
	ticker := time.NewTicker(5 * time.Second)

	go func() {
		for t := range ticker.C {
			fmt.Println("This code should run 5 secounds apart")
			fmt.Println("Tick at", t)
		}
	}()

	if err != nil {
		fmt.Printf("Push Notification Refresh Error: %v", err)
	}
	server.Idle()
}
