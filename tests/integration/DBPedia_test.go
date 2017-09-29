package tests

import (
	"testing"
	"nli-go/lib/common"
	"nli-go/lib/global"
)

func TestDBPedia(t *testing.T) {

	log := common.NewSystemLog(false)
	system := global.NewSystem(common.Dir()+"/../../resources/dbpedia/config-online.json", log)

	if !log.IsOk() {
		t.Errorf(log.String())
		return
	}

	var tests = []struct {
		question string
		answer   string
	}{
		{"Who married Lord Byron?", "Anne Isabella Milbanke married him"},
		{"Who married Anne Isabella Milbanke?", "George Gordon Byron married her"},
	}

	for _, test := range tests {

		log.Clear()

		answer := system.Answer(test.question)

		if answer != test.answer {
			t.Errorf("Test relationships: got %v, want %v", answer, test.answer)
			t.Error(log.String())
		}
	}
}