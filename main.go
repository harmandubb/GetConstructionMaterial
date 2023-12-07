package main

import (
	_ "docstruction/getconstructionmaterial/GCalls"
	server "docstruction/getconstructionmaterial/Server"
)

func main() {
	server.Idle()
}
