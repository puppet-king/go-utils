//go:build mage
// +build mage

package main

import (
	"go-utils/utils/registry"
	_ "go-utils/utils/strutil" // 导入触发 init 注册函数
)

type Args []string

func TrimHello() {
	registry.CmdExecute(Args{"strutil:Trim", "  hello  "})
}
