package binding

import (
	"testing"
	"tomm/log"
)

type TestStruct struct {
	A string `form:"a_field" default:"aaaa"`
	B string `form:"b_field" default:"bbbb"`
}

func TestDefaultBind(t *testing.T) {
	ts := TestStruct{}

	log.Info("TestDefaultBind Bind Name is %s", formBind.Name())
	testForm := make(map[string][]string)
	testForm["a_field"] = []string{""}
	testForm["b_field"] = []string{""}
	err := formBind.testInterface(testForm, &ts)
	if err != nil {
		log.Error("testInterface error %s ", err.Error())
		return
	}

	log.Info("TestStruct content is %v ", ts)
}
