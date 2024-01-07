package GPT

import (
	g "docstruction/getconstructionmaterial/GCalls"
	"encoding/json"
	"fmt"
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

func TestParseGPTAnalysisMaterialResponse(t *testing.T) {
	srv := g.ConnectToGmailAPI()

	user := "info@docstruction.com"

	unreadMsgData, err := g.GetUnreadMessagesData(srv, user)
	if err != nil {
		t.Error(err)
	}

	emailInfo, _, err := g.GetMessage(srv, unreadMsgData.Messages[0], user)
	if err != nil {
		t.Error(err)
	}

	templatePath := "./GPT_Prompts/email_material_check_prompt.txt"
	prompt, err := createReceiceEmailAnalysisPrompt(templatePath, emailInfo.Body)
	if err != nil {
		t.Error(err)
	}

	resp, err := promptGPT(prompt)
	if err != nil {
		t.Error(err)
	}

	emailProductInfo, err := parseGPTAnalysisMaterialResponse(resp)

	fmt.Println(emailProductInfo)

	res2B, _ := json.Marshal(emailProductInfo)
	fmt.Println(string(res2B))
}
