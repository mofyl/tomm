package binding

import (
	"go.uber.org/zap"
	"testing"
	"tomm/log"
)

type TestStruct struct {
	A string `form:"a_field" default:"aaaa"`
	B string `form:"b_field" default:"bbbb"`
}

func TestDefaultBind(t *testing.T) {
	ts := TestStruct{}

	log.Info("TestDefaultBind", zap.String("binding Name", formBind.Name()))
	testForm := make(map[string][]string)
	testForm["a_field"] = []string{""}
	testForm["b_field"] = []string{""}
	err := formBind.testInterface(testForm, &ts)
	if err != nil {
		log.Error("testInterface error ", zap.String("err", err.Error()))
		return
	}

	log.Info("TestStruct ", zap.Any("content", ts))
}
