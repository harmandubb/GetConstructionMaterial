package server

import (
	"fmt"
	"testing"
)

func TestPromptGPTMaterialCategorization(t *testing.T) {
	prompt, err := CreateMaterialCategorizationPrompt("./material_catigorization_prompt.txt", "Fire Stopping Collars")
	if err != nil {
		t.Fail()
	}

	resp, err := promptGPT(prompt)
	if err != nil {
		t.Fail()
	}

	fmt.Println(resp)

}
