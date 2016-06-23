package blog

// 常用的类型转换

import (
	"encoding/base64"
	"path"
	"strconv"
	"strings"
)

var (
	// 能被上传的文件扩展名
	IMAGES_EXT = []string{".jpg", ".jpeg", ".gif", ".png", ".bmp"}
)

const INVALID_VALUE = 0xFFFFFFFF

// 整型转字符串
func itos(i int) string {
	return strconv.Itoa(i)
}

// 字符串转整型
func stoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func stoui(s string) uint {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0xFFFFFFFF
	}
	return uint(i)
}

// 字符串转逻辑型
func stob(s string) bool {
	b, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}
	return b
}

// 逻辑型转字符串
func btos(b bool) string {
	return strconv.FormatBool(b)
}

// 检查上传的文件的护展名
func isCanUploadFile(file string) bool {
	for _, ext := range IMAGES_EXT {
		if path.Ext(strings.ToLower(file)) == ext {
			return true
		}
	}
	return false
}

// 获取文件扩展名
func getFileExt(file string) string {
	return strings.ToLower(path.Ext(file))
}

func base64EnString(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func base64DeString(str string) string {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return ""
	}
	return string(data)
}

func shortArticle(text string, length int) string {
	if length > 0 {
		if len(text) > length {
			return string(text[:length])
		}
	}
	return text
}

// 判断指定字符串是否在数组中
func stringInArray(arr []string, s string) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}
	return false
}

// 一个简单的html摘要生成
func simpleHTMLSummary(htmlstr string, length int) string {
	istag := false
	count := 0
	htmlindex := 0
	for i, s := range htmlstr {
		switch s {
		case '<':
			istag = true
		case '>':
			istag = false
		default:
			if !istag {
				count += 1
				// 收集到够长的字符了， 返回html索引位置
				if count > length {
					htmlindex = i
					goto endchar
				}
			}
		}
	}
endchar:

	if count <= length {
		return htmlstr
	}

	if htmlindex != 0 {
		// 获取指定字数的HTML，除html标记的纯文本字数计算出的
		shortstr := string(htmlstr[:htmlindex])
		// 提取打开标签
		opentags := []string{}
		j := -1
		for i, s := range shortstr {
			switch s {
			case '<':
				j = i
			case ' ', '>':
				if j != -1 {
					if shortstr[j+1] != '/' {
						opentags = append(opentags, string(shortstr[j+1:i]))
					}
					j = -1
				}
			}
		}
		// 提取闭合标签
		j = -1
		cloasetags := []string{}
		for i := len(shortstr) - 1; i >= 0; i-- {
			switch shortstr[i] {
			case '>':
				j = i
			case '<':
				if j != -1 {
					if shortstr[i+1] == '/' {
						cloasetags = append(cloasetags, string(shortstr[i+2:j]))
					}
					j = -1
				}
			}
		}
		if len(opentags) == len(cloasetags) {
			return shortstr
		}
		shortstr += " [......]"
		noclosetags := []string{"br", "hr", "input", "img"}
		for _, o := range opentags {
			if stringInArray(noclosetags, o) {
				continue
			}
			if !stringInArray(cloasetags, o) {
				shortstr += "</" + o + ">"
			}
		}
		return shortstr
	}
	return ""
}

// 从 xxx.xxx.xxx.xxx:xxx格式中取出ip地址
func getIpAddrByStr(ip string) string {
	i := strings.LastIndex(ip, ":")
	if i == -1 {
		return ip
	}
	return ip[:i]
}
