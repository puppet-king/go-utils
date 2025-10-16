// Package fileutil */
package fileutil

import (
	"fmt"
	"go-utils/utils/registry"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"time"
)

const PackageName = "fileutil"

func init() {
	registry.BatchRegister(PackageName, []interface{}{
		DirStat,
	})
}

// FileStat 文件信息
type FileStat struct {
	Name string
	Path string
	Size int64
	Type string
}

type Summary struct {
	TotalSize  int64
	TypeSize   map[string]int64
	TypeCount  map[string]int
	LargeFiles int
	LargeSize  int64
}

// DirStat 扫描目录，统计文件信息，按大小降序排序，展示表格，并可保存大文件
func DirStat(dir string, sizeThreshold int64, save bool) error {
	var files []FileStat
	summary := &Summary{
		TypeSize:  make(map[string]int64),
		TypeCount: make(map[string]int),
	}

	var outFile = ""
	if save {
		outFile = fmt.Sprintf("tmp/%s_file.txt", GenDateRandomFileName())
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		ftype := detectType(info.Name())
		files = append(files, FileStat{
			Name: info.Name(),
			Path: path,
			Size: info.Size(),
			Type: ftype,
		})

		// 更新汇总
		summary.TotalSize += info.Size()
		summary.TypeSize[ftype] += info.Size()
		summary.TypeCount[ftype]++

		if sizeThreshold > 0 && info.Size() >= sizeThreshold {
			summary.LargeFiles++
			summary.LargeSize += info.Size()
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 输出表格
	fmt.Printf("当前目录: %s 总大小: %s\n\n", dir, humanSize(summary.TotalSize))
	// 将 map 转成 slice 并排序
	typeSizeList := make([]struct {
		Type  string
		Size  int64
		Count int
	}, 0, len(summary.TypeSize))

	for t, sz := range summary.TypeSize {
		typeSizeList = append(typeSizeList, struct {
			Type  string
			Size  int64
			Count int
		}{t, sz, summary.TypeCount[t]})
	}

	// 按总大小降序排序
	sort.Slice(typeSizeList, func(i, j int) bool {
		return typeSizeList[i].Size > typeSizeList[j].Size
	})

	// 输出汇总表格
	fmt.Printf("%-10s %10s %10s\n", "类型", "大小", "占比")
	for _, ts := range typeSizeList {
		ratio := float64(ts.Size) / float64(summary.TotalSize) * 100
		ratioStr := fmt.Sprintf("%.1f%%", ratio)
		fmt.Printf("%-8s %16s %12s \n", ts.Type, humanSize(ts.Size), ratioStr)
	}

	// 保存大文件
	if sizeThreshold > 0 {
		fmt.Printf("\n大于 %s 的文件数量: %d, 总大小: %s\n",
			humanSize(sizeThreshold), summary.LargeFiles, humanSize(summary.LargeSize))

		var f *os.File
		var err error
		if save {
			f, err = os.Create(outFile)
			if err != nil {
				return err
			}
			defer f.Close()
		}

		fmt.Printf("%-8s %-30s %-10s\n", "类型", "名称", "大小")
		for _, file := range files {
			if file.Size >= sizeThreshold {
				if save {
					fmt.Fprintf(f, "%s %s %s\n", file.Type, file.Path, humanSize(file.Size))
				}

				fmt.Printf("%-8s %-30s %10s\n",
					file.Type, file.Name, humanSize(file.Size))

			}
		}

		fmt.Printf("\n")
		if save {
			absPath, _ := filepath.Abs(outFile)
			fmt.Println("大文件列表已保存至:", absPath)
		}
	}

	return nil
}

// detectType 根据文件扩展名识别类型
func detectType(name string) string {
	switch ext := filepath.Ext(name); ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".svg":
		return "image"
	case ".js":
		return "js"
	case ".css":
		return "css"
	case ".gz":
		return "gzip"
	default:
		return "other"
	}
}

// humanSize 转换字节为可读格式
func humanSize(size int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	switch {
	case size >= GB:
		return fmt.Sprintf("%.2fGB", float64(size)/GB)
	case size >= MB:
		return fmt.Sprintf("%.2fMB", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.2fKB", float64(size)/KB)
	default:
		return fmt.Sprintf("%dB", size)
	}
}

// GenDateRandomFileName 生成文件名标识：YYYYMMDD + 4位随机数
func GenDateRandomFileName() string {
	now := time.Now()
	date := now.Format("20060102") // YYYYMMDD

	// 随机数
	rand.Seed(time.Now().UnixNano())
	rnd := rand.Intn(10000) // 0~9999
	return fmt.Sprintf("%s%04d", date, rnd)
}
