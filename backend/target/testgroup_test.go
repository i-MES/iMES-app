package target

import (
	"testing"

	"github.com/i-mes/imes-app/backend/utils"
)

func TestParsePython(t *testing.T) {
	fp := new(Parser)
	ap := utils.GetAppPath()
	err := fp.ParsePython(1, ap+"/testcase/python/test_led.py")
	if err != nil {
		t.Error(err)
	}
}
