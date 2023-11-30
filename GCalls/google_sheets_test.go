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
