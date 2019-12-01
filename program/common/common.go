package common

import (
	"os"
	"sort"
	"strings"
	"unicode"
)

/* 公共小函数库 */

// PathExists 判断文件或目录是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// StrFirstToUpper 字符串首字母转大些
func StrFirstToUpper(str string, isFirst bool) string {
	if str == "" {
		return ""
	}
	temp := strings.Split(str, "_")
	for i := 0; i < len(temp); i++ {
		if isFirst == false && i == 0 {
			temp[i] = FirstToLower(temp[i])
			continue
		}
		if IsStartUpper(temp[i]) == false {
			temp[i] = FirstToUpper(temp[i])
		}
	}

	return strings.Join(temp, "")
}

// IsStartUpper 判断首字母是否是大写
func IsStartUpper(s string) bool {
	return unicode.IsUpper([]rune(s)[0])
}

// FirstToUpper 首字母转大写
func FirstToUpper(str string) string {
	var upperStr string
	vv := []rune(str)
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 97 && vv[i] <= 122 {
				vv[i] -= 32
				upperStr += string(vv[i])
			} else {
				// fmt.Println("Not begins with lowercase letter,")
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}

// FirstToLower 首字母转小写
func FirstToLower(str string) string {
	var upperStr string
	vv := []rune(str)
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 65 && vv[i] <= 90 {
				vv[i] += 32
				upperStr += string(vv[i])
			} else {
				// fmt.Println("Not begins with lowercase letter,")
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}

// LsPath 获取文件列表 - 不包含 .DS_Store Desktop.ini
func LsPath(filePath string) (pathList []string, err error) {
	f, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer f.Close()
	pathList, err = f.Readdirnames(-1)
	for k, v := range pathList {
		if v == "Desktop.ini" || v == ".DS_Store" {
			pathList = append(pathList[:k], pathList[k+1:]...)
		}
	}
	sort.Strings(pathList)
	return
}
