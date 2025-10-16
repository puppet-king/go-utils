package strutil

import (
	"go-utils/utils/registry"
	"strings"
)

const PackageName = "strutil"

// 包初始化时自动注册函数
func init() {
	registry.BatchRegister(PackageName, []interface{}{
		Trim,
		ToUpper,
	})
}

// Trim 修剪字符串两端的空白（包括全角空格）
func Trim(s string) string {
	return strings.Trim(s, " \t\n\r　")
}

// ToUpper 转换为大写
func ToUpper(s string) string {
	return strings.ToUpper(s)
}

// toUpper 转换为大写
func toUpper(s string) string {
	return strings.ToUpper(s)
}
