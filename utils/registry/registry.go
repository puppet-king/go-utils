package registry

import (
	"fmt"
	"github.com/fatih/color"
	"reflect"
	"runtime"
	"strings"
	"sync"
)

// FunctionInfo 存储函数信息
type FunctionInfo struct {
	Package string
	Name    string
	Func    reflect.Value
}

var (
	funcMap = make(map[string]FunctionInfo)
	mu      sync.Mutex
)

// register 注册函数（由工具包中的init函数调用）
func register(pkgName string, fn interface{}) bool {
	funcName, ok := functionName(fn)
	if !ok {
		return false
	}

	mu.Lock()
	defer mu.Unlock()
	key := pkgName + ":" + funcName
	funcMap[key] = FunctionInfo{
		Package: pkgName,
		Name:    funcName,
		Func:    reflect.ValueOf(fn),
	}

	return true
}

// BatchRegister 批量注册
func BatchRegister(pkgName string, fn []interface{}) {
	var total = len(fn)
	var count = 0

	for _, f := range fn {
		if register(pkgName, f) {
			count++
		}
	}

	if total == count {
		color.Green("%s register success \n", pkgName)
	} else {
		color.Red("%s register failed \n", pkgName)
	}
}

func functionName(fn interface{}) (string, bool) {
	val := reflect.ValueOf(fn)
	if val.Kind() != reflect.Func {
		return "", false
	}

	pc := val.Pointer()
	fnInfo := runtime.FuncForPC(pc)
	fullName := fnInfo.Name()
	parts := strings.Split(fullName, ".")
	funcName := parts[len(parts)-1]

	// 如果函数名是小写开头（未导出），跳过注册
	if len(funcName) == 0 || funcName[0] >= 'a' && funcName[0] <= 'z' {
		return "", false
	}

	return funcName, true
}

// Function 根据包名和函数名获取函数
func Function(pkgName, funcName string) (reflect.Value, bool) {
	key := pkgName + ":" + funcName
	fn, ok := funcMap[key]
	return fn.Func, ok
}

// AllFunctions 获取所有已注册的函数
func AllFunctions() []FunctionInfo {
	var funcs []FunctionInfo
	for _, fn := range funcMap {
		funcs = append(funcs, fn)
	}
	return funcs
}

// Execute 执行命令
func Execute(pkgName, funcName string, parts []string) bool {
	// 动态查找函数
	fn, exists := Function(pkgName, funcName)
	if !exists {
		fmt.Printf("未找到函数: %s:%s\n", pkgName, funcName)
		return false
	}

	// 转换参数
	args := parts[1:]
	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}

	// 检查参数数量是否匹配
	if fn.Type().NumIn() != len(in) {
		fmt.Printf("参数数量错误，需要 %d 个参数，实际提供 %d 个\n",
			fn.Type().NumIn(), len(in))
		return false
	}

	// 调用函数
	results := fn.Call(in)

	// 处理返回结果
	for _, result := range results {
		if result.CanInterface() {
			if err, ok := result.Interface().(error); ok {
				fmt.Printf("错误: %v\n", err)
				return false
			}
			fmt.Println(result.Interface())
		}
	}

	return true
}

// CmdExecute mage 执行命令
func CmdExecute(args []string) bool {
	if len(args) == 0 {
		color.Red("命令格式错误，正确格式: 包名:函数名 [参数...]")
		return false
	}

	cmdParts := strings.Split(args[0], ":")
	if len(cmdParts) != 2 {
		color.Red("命令格式错误，正确格式: 包名:函数名 [参数...]")
		return false
	}

	pkgName, funcName := cmdParts[0], cmdParts[1]

	return Execute(pkgName, funcName, args)
}
