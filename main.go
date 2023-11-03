package main

import (
	"docstruction/getconstructionmaterial/api"
)

func main() {
	// api.SendEmail()
	// api.CheckDataBase()
	api.AddProductBasic("Meta Caulk Collars", "Fire Stopping", 10.01)
	api.AddProductPicture("Meta Caulk Collars", "./images/img1.jpg")

}
