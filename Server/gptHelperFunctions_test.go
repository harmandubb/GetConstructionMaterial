package server

import (
	"testing"
)

func TestPromptGPTMaterialCategorization(t *testing.T) {
	category, err := PromptGPTMaterialCatogorization("./material_catigorization_prompt.txt", "bx 4 connector wire")
	if err != nil {
		t.Error(err)
	}

	if category != "Electrical" {
		t.Fail()
	}

}
