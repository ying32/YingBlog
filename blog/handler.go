package blog

// 此文件是参考gopher来写的，所以文件名和函数几乎一样

// 用作渲染用

import (
	//	"database/sql"
	"encoding/json"
	"html/template"
	"net/http"
	texttemplate "text/template"
	"time"

	"github.com/gorilla/mux"
)

// 回调的参数类型
type Handlder struct {
	http.ResponseWriter
	*http.Request
	HTMLFile  string // 渲染的静态文件
	StartTime time.Time
	DB        *MDB //*sql.DB
}

// 渲染文件
func (self *Handlder) render(file string, datas ...map[string]interface{}) {
	var data = make(map[string]interface{})
	if len(datas) == 1 {
		data = datas[0]
	} else if len(datas) != 0 {
		panic("数据传入过多")
	}
	// 应该放到template
	tpl, err := template.ParseFiles("./template/" + file)
	if err != nil {
		panic(err)
	}
	self.setHeaderResponseType("text/html")
	// 默认数据
	self.setDefaultData(data)
	// 参数要有默认的东东，比如footer，
	tpl.Execute(self.ResponseWriter, data)
}

// 默认数据
func (self *Handlder) setDefaultData(data map[string]interface{}) {
	data["goVersion"] = goVersion
	data["startTime"] = self.StartTime
	data["website"] = appSvrConf.WebSite
}

func (self *Handlder) renderTemplate(file, baseFile string, datas ...map[string]interface{}) {
	var data = make(map[string]interface{})
	if len(datas) == 1 {
		data = datas[0]
	} else if len(datas) != 0 {
		panic("数据传入过多")
	}
	// 默认数据
	self.setDefaultData(data)
	self.ResponseWriter.Write(parseTemplate(file, baseFile, data))
}

// 返回xml文件
func (self *Handlder) renderXML(file string, datas ...map[string]interface{}) {
	var data = make(map[string]interface{})
	if len(datas) == 1 {
		data = datas[0]
	} else if len(datas) != 0 {
		panic("数据传入过多")
	}
	// 应该放到template
	tpl, err := texttemplate.ParseFiles("./template/" + file)
	if err != nil {
		panic(err)
	}
	//	self.setDefaultData(data)
	data["website"] = appSvrConf.WebSite
	self.setHeaderResponseType("application/xml")
	tpl.Execute(self.ResponseWriter, data)
}

// 参数提取
func (self *Handlder) param(name string) string {
	return mux.Vars(self.Request)[name]
}

// 设置头数据
func (self *Handlder) setHeaderResponseType(datatype string) {
	self.ResponseWriter.Header().Set("Content-Type", datatype+"; charset=utf-8")
	//self.ResponseWriter.Header().Set("Content-Type", datatype)
	//self.ResponseWriter.Header().Set("charset", "utf-8")
}

// 重定向
func (self *Handlder) redirect(url string, statecode int) {
	http.Redirect(self.ResponseWriter, self.Request, url, statecode)
}

// 返回json数据
func (self *Handlder) renderJson(data interface{}) {
	jsbytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	self.setHeaderResponseType("text/json") //ie下使用application/json竟然会有问题，唉
	self.ResponseWriter.Write(jsbytes)
}

// 返回纯文本数据
func (self *Handlder) renderText(text string) {
	self.ResponseWriter.Write([]byte(text))
}

// Q开头表示 Query 也是就是Form中的
// P开头表示 Post表单中的

// Query表单中的
func (self *Handlder) QInt(key string) int {
	return stoi(self.FormValue(key))
}

func (self *Handlder) QStr(key string) string {
	return self.FormValue(key)
}

func (self *Handlder) QBool(key string) bool {
	return stob(self.FormValue(key))
}

func (self *Handlder) QEqStr(key, value string) bool {
	return self.QStr(key) == value
}

func (self *Handlder) QEqInt(key string, value int) bool {
	return self.QInt(key) == value
}

// Post表单中的
func (self *Handlder) PInt(key string) int {
	return stoi(self.PostFormValue(key))
}

func (self *Handlder) PUint(key string) uint {
	return stoui(self.PostFormValue(key))
}

func (self *Handlder) PEqStr(key, value string) bool {
	return self.PStr(key) == value
}

func (self *Handlder) PStr(key string) string {
	return self.PostFormValue(key)
}

func (self *Handlder) PBool(key string) bool {
	return stob(self.PostFormValue(key))
}

func (self *Handlder) PEqInt(key string, value int) bool {
	return self.PInt(key) == value
}

func (self *Handlder) PEqUInt(key string, value uint) bool {
	return self.PUint(key) == value
}

func (self *Handlder) redirectError() {
	self.redirect("/error", http.StatusFound)
}
