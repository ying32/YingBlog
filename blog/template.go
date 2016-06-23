package blog

// 改自gopher template.go单元

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"time"
)

// 两个基础模版
const (
	BASEHTML  = "base.html"
	ADMINHTML = "admin/base.html"
)

func loadtimes(startTime time.Time) string {
	return fmt.Sprintf("%dms", time.Now().Sub(startTime)/1000000)
}

func html(text string) template.HTML {
	return template.HTML(text)
}

func shortarticle(text string, length int) template.HTML {
	if length > 0 {
		if len(text) > length {
			return template.HTML(text[:length])
		}
	}
	return template.HTML(text)
}

func formatdatetime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// 模版中自定义的解析函数
var templateFuncMaps = template.FuncMap{
	"html": html,
	"include": func(filename string, data map[string]interface{}) template.HTML {
		var buf bytes.Buffer
		t := template.New(filename).Funcs(template.FuncMap{"html": html, "loadtimes": loadtimes})
		t, err := t.ParseFiles("./template/" + filename)
		if err != nil {
			panic(err)
		}
		err = t.Execute(&buf, data)
		if err != nil {
			panic(err)
		}
		return template.HTML(buf.Bytes())
	},
	// 纯复制
	"loadtimes":      loadtimes,
	"shortarticle":   shortarticle,
	"base64decode":   base64DeString,
	"formatdatetime": formatdatetime}

// 分析模版
func parseTemplate(file, baseFile string, data map[string]interface{}) []byte {
	var buf bytes.Buffer
	t := template.New(file).Funcs(templateFuncMaps)
	basebytes, err := ioutil.ReadFile("./template/" + baseFile)
	if err != nil {
		panic(err)
	}
	t, err = t.Parse(string(basebytes))
	if err != nil {
		panic(err)
	}
	t, err = t.ParseFiles("./template/" + file)
	if err != nil {
		panic(err)
	}
	err = t.Execute(&buf, data)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func renderTemplate(handlder *Handlder, file, baseFile string, data map[string]interface{}) {
	data["goVersion"] = goVersion
	handlder.ResponseWriter.Write(parseTemplate(file, baseFile, data))
}
