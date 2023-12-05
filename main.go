package main

import (
	_ "docstruction/getconstructionmaterial/GCalls"
	server "docstruction/getconstructionmaterial/Server"
	"fmt"
)

func main() {
	fmt.Println("Running the Main function")
	server.Idle()
}
