package gcalls

import (
	"testing"
	"time"
)

func TestPublish(t *testing.T) {
	success := sendEmailInfo(time.Now(), "harmandubb@docstruction.com", "1ZowyzJ008toPYNn0mFc2wG6YTAop9HfnbMPLIM4rRZw")

	if !success {
		t.Fail()
	}
}
