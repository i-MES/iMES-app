package python

/*
#cgo pkg-config: python3-embed
#define PY_SSIZE_T_CLEAN
#include <Python.h>
*/
import "C"

import (
	"fmt"
	"regexp"
)

type Info struct {
	version string
}

func Version_Check(ver string) bool {
	_ver := Py_GetVersion()
	matched, err := regexp.Match(ver, []byte(_ver))
	if err == nil {
		return matched
	}
	return false
}

func PyErr_Occurred() bool {
	if err := C.PyErr_Occurred(); togo(err) != nil {
		return true
	} else {
		return false
	}
}

// 打印 python 中的 stderr
func PyErr_Print() string {
	errstr := ""
	if err := C.PyErr_Occurred(); togo(err) != nil {
		fmt.Println("XXXXXX python exception XXXXXX")
		// sys.stdout = io.StringIO() // 修改 stdout 为 StringIO 实例
		// sys.stderr = io.StringIO()
		C.PyErr_Print()
		// errstr := sys.stderr.getvalue() // 读出 stdout 的错误打印
		// sys.stdout.close()
		// sys.stderr.close()
		fmt.Println("XXXXXX end XXXXXX")
	}
	return errstr
}

func InitLog() {
	py := `
from logging import getLogger, getLevelName, Formatter, StreamHandler

log = getLogger()
log.setLevel(getLevelName('INFO'))
log_formatter = Formatter("%(asctime)s [%(process)d] [%(thread)d] [%(levelname)s] %(name)s: %(message)s")

console_handler = StreamHandler()
console_handler.setFormatter(log_formatter)
log.addHandler(console_handler)

log.info("Set Python Log")
`
	PyRun_SimpleString(py)
}

func LogInfo(msg string) {
	PyRun_SimpleString(fmt.Sprintf("getLogger().info('%s: ')", msg))
}

func LogProcessId() {
	PyRun_SimpleString("print('Python process id: ', os.getpid())")
}

func LogThreadId() {
	PyRun_SimpleString("print('Python threading id: ', threading.current_thread().native_id)")
}
