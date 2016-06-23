package blog

// 此单元是参照gopher写的，所以文件名和函数名几乎一样

import (
	"fmt"
	"net/http"
	"time"
)

// 回调函数
type HandleFunc func(*Handlder)

// 路由结构体
type Route struct {
	URL        string     // 请求的url地址
	HTMLFile   string     // 渲染的静态文件
	Permission PerType    // 权限
	HandleFunc HandleFunc // 回调函数
}

// 新建带有请求的handler,在这里可以启动数据库等等操作吧
func NewHandler(w http.ResponseWriter, r *http.Request, file string) *Handlder {
	return &Handlder{
		ResponseWriter: w,
		Request:        r,
		StartTime:      time.Now(),
		HTMLFile:       file,
		DB:             openDatabase(), // 每次都去打开数据库，不知道好不好，还是一次性打开呢？？？？
	}
}

// 纯复制的。。。
func fileHandler(w http.ResponseWriter, req *http.Request) {
	url := req.Method + " " + req.URL.Path
	fmt.Println(url)
	filePath := req.URL.Path[1:]
	//http.
	http.ServeFile(w, req, filePath)
}

// 自定义404页面
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	// why?
	if r.RequestURI == "/favicon.ico" {
		return
	}
	http.Redirect(w, r, "/error", http.StatusFound)
}

// 重定向
func (self *Handlder) Redirect(url string) {
	http.Redirect(self.ResponseWriter, self.Request, url, http.StatusFound)
}

// 路由规则表
// 配置种各样的路由规则
var routes = []Route{
	{"/", "/blog.html", Everyone, indexHandler},
	{"/index", "blog.html", Everyone, indexHandler},
	{"/Me", "aboutMe.html", Everyone, aboutMeHandler},
	{"/rss", "rss.xml", Everyone, rssHandler},
	// editor
	//{"/editor", "editor.html", Administrator, editorHandler},
	//{"/robots.txt", "robots.txt", Everyone, robotsHandler},
	// admin
	{"/admin/login", "login.html", Everyone, loginHandler},

	{"/error", "404.html", Everyone, errorHandler},

	{"/admin", "admin/index.html", Administrator, adminHandler},
	{"/admin/index", "admin/index.html", Administrator, adminHandler},
	{"/admin/user", "admin/user.html", Administrator, adminHandler},
	{"/admin/log", "admin/log.html", Administrator, adminHandler},
	{"/admin/articlemgr", "admin/articlemgr.html", Administrator, articleMgrHandler},
	{"/admin/form", "admin/form.html", Administrator, adminHandler},
	{"/admin/gallery", "admin/gallery.html", Administrator, adminHandler},
	{"/admin/logout", "", Administrator, logoutHandler},

	{"/admin/upload", "", Administrator, uploadHandler},
	{"/admin/savearticle", "", Administrator, saveArticleHandler},
	{"/admin/messages", "admin/messages.html", Administrator, messagesHandler},
	{"/admin/statistic", "/admin/statistic.html", Administrator, statisticHandler},

	{"/article/{id}", "article.html", Everyone, articleHandler},
	{"/admin/article/del", "", Administrator, articleDeleteHandler},
	{"/admin/article/edit/{id}", "editor.html", Administrator, articleEditHandler},
	{"/admin/article/new", "editor.html", Administrator, articleNewHandler},
	{"/admin/article/operation", "", Administrator, articleOperationHandler},
}
