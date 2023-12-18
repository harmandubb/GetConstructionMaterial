package gcalls

import (
	"testing"
	"time"
)

func TestPublis(t *testing.T) {
	success := SendEmailInfo(time.Now(), "harmandubb@docstruction.com", "1ZowyzJ008toPYNn0mFc2wG6YTAop9HfnbMPLIM4rRZw")

	if !success {
		t.Fail()
	}
}

func TestMaterialFormPublish(t *testing.T) {
	matInfo := MaterialFormInfo{
		Email:    "harmandubb@docstruction.com",
		Material: "Fire Stopping Collars",
		Location: "Surrey BC",
	}

	success := SendMaterialFormInfo("1NXTK2G6sQOs0ZSQ1046ijoanPDNWPKOc0-I7dEMotQ8", matInfo)

	if !success {
		t.Fail()
	}
}
