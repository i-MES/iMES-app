package testset

import (
	"fmt"
	"runtime"

	"github.com/i-mes/imes-app/backend/python"
	"github.com/i-mes/imes-app/backend/utils"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// 测试项
type TestItem struct {
	Title    string `json:"title"`
	Desc     string `json:"desc"`
	FileName string `json:"filename"`
	FuncName string `json:"funcname"`
	Sequence int    `json:"sequence"`
}

func InitTestItems() {
	data := make([]TestItem, 0)
	data = append(data,
		TestItem{"MCU Test", "MCU Test...", "", "test_mcu", 1},
		TestItem{"Memory Test", "Memory Test...", "", "test_memory", 2},
		TestItem{"Network Test", "Network Test...", "", "test_network", 3},
	)
	SaveTestItems(data)
}

func SaveTestItems(data []TestItem) {
	_data := make(map[string]interface{})
	_data["testitem"] = data
	utils.OutputConfigData(_data)
}
func LoadTestItems() []TestItem {
	var data []TestItem
	_data := []byte(json.Get(utils.InputConfigData("testitem")).ToString())
	err := json.Unmarshal(_data, &data)
	if err != nil {
		fmt.Println(err)
	}
	return data
}

func (ti *TestItem) Run(tg_name string) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	fmt.Println(ti)

	if !python.Py_IsInitialized() {
		python.Py_Initialize()
	}
	fmt.Println("-=")
	_gil := python.PyGILState_Ensure()
	defer python.PyGILState_Release(_gil)

	fmt.Println("-=-=")
	_mod := python.PyImport_ImportFile(ti.FileName)
	defer _mod.DecRef()
	if _mod == nil {
		panic("mod error")
	}

	pFunc := _mod.GetAttrString(ti.FuncName)
	if pFunc != nil {
		pValue := pFunc.CallObject(nil)
		if pValue != nil {
			fmt.Println(pValue)
		}
	}
	fmt.Println("Run TI: ", ti.Title, ti.Desc, ti.FuncName, ti.Sequence)
}
