package gcalls

import (
	"fmt"
	"testing"
)

func TestConnectToGmail(t *testing.T) {
	srv := ConnectToGmail()

	fmt.Println(srv)
}
