package target

import (
	"fmt"

	"github.com/i-mes/imes-app/backend/utils"
)

// 测试项
type TestItem struct {
	Id    string   `json:"id"`    // uuid
	Title string   `json:"title"` // function name
	Desc  string   `json:"desc"`  // python docstr
	Args  []string `json:"args"`
}

// func CreateTestItemExample() {
// 	data := make([]TestItem, 0)
// 	data = append(data,
// 		TestItem{utils.UUID(), "MCU Test", "MCU Test...", "", "test_mcu", 1},
// 		TestItem{utils.UUID(), "Memory Test", "Memory Test...", "", "test_memory", 2},
// 		TestItem{utils.UUID(), "Network Test", "Network Test...", "", "test_network", 3},
// 	)
// 	SaveTestItems(data)
// }

func SaveTestItems(data []TestItem) {
	_data := make(map[string]interface{})
	_data["testitem"] = data
	utils.OutputConfigData(_data)
}

func LoadTestItems() []TestItem {
	if data, err := utils.InputConfigData("testitem"); err == nil {
		var ti []TestItem
		_ti := []byte(json.Get(data).ToString())
		err := json.Unmarshal(_ti, &ti)
		if err != nil {
			fmt.Println(err)
		}
		return ti
	}
	return nil
}

// func (ti *TestItem) Run(tg_name string) {
// 	runtime.LockOSThread()
// 	defer runtime.UnlockOSThread()
// 	fmt.Println(ti)

// 	if !python.Py_IsInitialized() {
// 		python.Py_Initialize()
// 	}
// 	fmt.Println("-=")
// 	_gil := python.PyGILState_Ensure()
// 	defer python.PyGILState_Release(_gil)

// 	fmt.Println("-=-=")
// 	_mod := python.PyImport_ImportFile(ti.FileName)
// 	defer _mod.DecRef()
// 	if _mod == nil {
// 		panic("mod error")
// 	}

// 	pFunc := _mod.GetAttrString(ti.FuncName)
// 	if pFunc != nil {
// 		pValue := pFunc.CallObject(nil)
// 		if pValue != nil {
// 			fmt.Println(pValue)
// 		}
// 	}
// 	fmt.Println("Run TI: ", ti.Title, ti.Desc, ti.FuncName, ti.Sequence)
// }

type TestItemStatus struct {
	TestEntityId string `json:"testentityid"`
	TestGroupId  string `json:"testgroupid"`
	TestClassId  string `json:"testclassid"`
	TestItemId   string `json:"testitemid"`
	Status       string `json:"status"`
	TimeStamp    int64  `json:"timestamp"`
}
