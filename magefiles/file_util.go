//go:build mage
// +build mage

package main

import (
	"go-utils/utils/fileutil"
	_ "go-utils/utils/fileutil" // 触发注册
)

// 扫描 assets 目录示例
func ScanAssets() error {
	dir := "C:\\project\\chaoge\\chaoge-admin\\dist\\assets"
	threshold := int64(1024 * 1024) // 1MB

	// fmt.Printf("扫描目录: %s, 大小阈值: %dB\n", dir, threshold)
	return fileutil.DirStat(dir, threshold, true)
}
