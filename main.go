package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"go-utils/utils/registry"  // 函数注册器
	_ "go-utils/utils/strutil" // 导入工具包
)

func main() {
	fmt.Println("===== go-utils 交互式工具 =====")
	fmt.Println("输入 'help' 查看帮助，'exit' 退出程序")
	fmt.Println("命令格式: 包名:函数名 [参数...] (例如: strutil:Trim '  hello  ')")
	fmt.Println("==============================")

	// 创建标准输入扫描器
	scanner := bufio.NewScanner(os.Stdin)

	// 进入交互循环
	for {
		fmt.Print("> ") // 命令提示符
		if !scanner.Scan() {
			break // 输入结束时退出
		}

		// 获取并处理输入
		input := scanner.Text()
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		// 处理特殊命令
		switch input {
		case "exit", "quit":
			fmt.Println("程序退出")
			return
		case "help":
			printHelp()
			continue
		case "list":
			printAllFunctions()
			continue
		}

		// 解析并执行命令
		executeCommand(input)
	}

	// 检查扫描错误
	if err := scanner.Err(); err != nil {
		fmt.Printf("输入错误: %v\n", err)
	}
}

// 执行命令
func executeCommand(input string) {
	// 分割输入为命令和参数（支持带引号的参数）
	parts, err := splitCommand(input)
	if err != nil {
		fmt.Printf("解析错误: %v\n", err)
		return
	}

	if len(parts) < 1 {
		return
	}

	// 解析命令（格式：包名:函数名）
	cmd := parts[0]
	cmdParts := strings.Split(cmd, ":")
	if len(cmdParts) != 2 {
		fmt.Println("命令格式错误，正确格式: 包名:函数名 [参数...]")
		return
	}

	pkgName, funcName := cmdParts[0], cmdParts[1]
	registry.Execute(pkgName, funcName, parts)
}

// 打印帮助信息
func printHelp() {
	fmt.Println("帮助信息:")
	fmt.Println("  命令格式: 包名:函数名 [参数...]")
	fmt.Println("  特殊命令:")
	fmt.Println("    help   - 显示帮助信息")
	fmt.Println("    list   - 显示所有可用函数")
	fmt.Println("    exit   - 退出程序")
	fmt.Println("  示例:")
	fmt.Println("    strutil:Trim '  hello world  '")
	fmt.Println("    strutil:ToUpper 'hello'")
}

// 打印所有可用函数
func printAllFunctions() {
	fmt.Println("可用函数列表:")
	for _, fn := range registry.AllFunctions() {
		fmt.Printf("  %s:%s\n", fn.Package, fn.Name)
	}
}

// 分割命令行（支持带引号的参数）
func splitCommand(input string) ([]string, error) {
	var parts []string
	var current strings.Builder
	inQuotes := false
	quoteChar := byte(0)

	for i := 0; i < len(input); i++ {
		c := input[i]

		// 处理引号
		if c == '"' || c == '\'' {
			if inQuotes && c == quoteChar {
				inQuotes = false
				quoteChar = 0
			} else if !inQuotes {
				inQuotes = true
				quoteChar = c
			} else {
				current.WriteByte(c)
			}
			continue
		}

		// 处理空格（引号内的空格不分割）
		if c == ' ' && !inQuotes {
			if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
			continue
		}

		current.WriteByte(c)
	}

	// 检查未闭合的引号
	if inQuotes {
		return nil, fmt.Errorf("未闭合的引号")
	}

	// 添加最后一个参数
	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts, nil
}
