package server

import (
	_ "embed"
	"testing"
)

func TestIdle(t *testing.T) {
	Idle()
}

func TestClientTest(t *testing.T) {
	clientTest()
}

// //go:embed GPT_Prompts/material_catigorization_prompt.txt
// var catigorizationTemplate string

// //go:embed GPT_Prompts/email_prompt.txt
// var emailTemplate string
