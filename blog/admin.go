package blog

import (
	"fmt"
	"math"
	"net/http"
)

// 用户信息
type User struct {
	id            int
	usr           string
	pwd           string
	nickname      string
	isenabled     bool
	permission    int
	email         string
	qqnum         int
	address       string
	headimg       string
	signature     string
	lastlogintime int64
	lastloginip   string
}

//	"github.com/gorilla/mux"

// 重定向到登录
func redirectLogin(handler *Handlder) {
	handler.redirect("/admin/login", http.StatusFound)
}

// 当前用户
func currentUser(handler *Handlder) (*User, bool) {
	session, _ := sessionStore.Get(handler.Request, "session")
	username, ok := session.Values["username"]
	if !ok {
		return nil, false
	}
	rows, err := handler.DB.Query("select * from users where usr=?", username.(string))
	if err != nil {
		return nil, false
	}
	defer rows.Close()
	usr := &User{}
	if rows.Next() {
		rows.Scan(&usr.id, &usr.usr, &usr.pwd, &usr.nickname, &usr.isenabled,
			&usr.permission, &usr.email, &usr.qqnum, &usr.address, &usr.headimg,
			&usr.signature, &usr.lastloginip, &usr.lastloginip)
		return usr, true
	}
	return nil, false
}

// 检测当前用户是否已经登录了, 貌似过于简单了
func checkIsLogin(handler *Handlder) bool {
	if _, ok := currentUser(handler); ok {
		return true
	}
	return false
}

// admin页面
func adminHandler(handler *Handlder) {
	handler.renderTemplate(handler.HTMLFile, ADMINHTML, map[string]interface{}{})
}

// 登录操作
func loginHandler(handler *Handlder) {
	if handler.Request.Method == "POST" {
		if usr, ok := queryCurrentUser(handler, handler.PStr("username")); ok {
			// 帐号被禁用
			msg := &TMessage{}
			msg.Success = true
			if !usr.isenabled {
				//redirectLogin(handler)
				fmt.Println("帐号被禁用")
				msg.Success = false
				msg.Message = "帐号已停用"
				return
			}
			// 密码错误
			if !handler.PEqStr("password", usr.pwd) {
				fmt.Println("密码错误")
				//redirectLogin(handler)
				msg.Success = false
				msg.Message = "密码错啦，亲！"
				return
			}
			// 保存session
			session, _ := sessionStore.Get(handler.Request, "session")
			session.Values["id"] = usr.id
			session.Values["username"] = usr.usr
			session.Save(handler.Request, handler.ResponseWriter)

			msg.Message = "登录成功！"

			handler.renderJson(msg)
			// 这里倒时候不用这么做了，用ajax异步请求来提交
			//handler.redirect("/admin", http.StatusFound)
		} else {
			msg := &TMessage{false, "未知错误，可能用户不存在！"}
			handler.renderJson(msg)
		}
	} else {
		// 直接输入登录页面时检测用户登录状态
		if checkIsLogin(handler) {
			handler.redirect("/admin", http.StatusFound)
			return
		}
		handler.renderTemplate(handler.HTMLFile, BASEHTML, map[string]interface{}{
			"isSignPage": true})
	}
}

// 登出操作
func logoutHandler(handler *Handlder) {
	session, _ := sessionStore.Get(handler.Request, "session")
	for key, _ := range session.Values {
		delete(session.Values, key)
	}
	session.Save(handler.Request, handler.ResponseWriter)
	handler.redirect("/admin/login", http.StatusFound)
}

// 输出文章信息列表
func articleMgrHandler(handler *Handlder) {
	handler.renderTemplate(handler.HTMLFile, ADMINHTML, map[string]interface{}{
		"articles":  queryArticlesByPage(handler, 0, 15),
		"categroys": queryCategroys(handler)})
}

func articleDeleteHandler(handler *Handlder) {
	result := &TMessage{}
	id := handler.PInt("id")
	if id < 0 {
		result.Success = false
		result.Message = "删除失败, id错误"
		return
	}
	result.Success = deleteArticleById(handler, id)
	if !result.Success {
		result.Message = "删除失败！"
	} else {
		result.Message = "删除成功！"
	}
	handler.renderJson(result)
}

func messagesHandler(handler *Handlder) {
	handler.renderTemplate(handler.HTMLFile, ADMINHTML, map[string]interface{}{})
}

func saveArticleHandler(handler *Handlder) {
	if handler.Request.Method == "POST" {
		msg := &TArticleMsg{}
		msg.InsertId = -1
		if handler.PEqInt("categroy", -1) {
			msg.Success = false
			msg.Message = "文章分类错误"
			return
		}
		if !handler.PEqStr("title", "") {
			// 判断是否有id

			if !handler.PEqUInt("id", INVALID_VALUE) {
				msg.Success = updateArticle(handler)
				if msg.Success {
					msg.Message = "更新文章成功"
				} else {
					msg.Message = "更新文章失败"
				}
			} else {
				msg.InsertId, msg.Success = insertNewArticle(handler)
				if msg.Success {
					msg.Message = "添加文章成功"
				} else {
					msg.Message = "添加文章失败"
				}
			}
		} else {
			msg.Success = false
			msg.Message = "添加失败，标题不能为空！"
		}
		//
		//if handler.p
		// remember
		//		cookie, err := handler.Cookie("login")
		//		if err == nil {
		//cookie.
		//		}

		handler.renderJson(msg)
	} else {
		handler.redirect("", http.StatusBadRequest)
	}
}

func articleEditHandler(handler *Handlder) {
	handler.renderTemplate(handler.HTMLFile, BASEHTML, map[string]interface{}{
		"article":   queryArticleById(handler, stoi(handler.param("id"))),
		"categroys": queryCategroys(handler),
		"isedit":    true,
		"iseditor":  true})
}

func articleNewHandler(handler *Handlder) {
	handler.renderTemplate(handler.HTMLFile, BASEHTML, map[string]interface{}{
		"categroys": queryCategroys(handler),
		"isedit":    false,
		"iseditor":  true})
}

func articleOperationHandler(handler *Handlder) {
	if handler.Request.Method == "POST" {
		msg := &TMessage{}
		msg.Success = setArticlePublic(handler)
		if msg.Success {
			if handler.PBool("ispublic") {
				msg.Message = "已公开"
			} else {
				msg.Message = "已隐藏"
			}
		} else {
			msg.Message = "更新失败"
		}
		handler.renderJson(msg)
	}
}

//---------------统计相关------------------------

//    <li class="am-disabled"><a href="#">«</a></li>
//      <li class="am-active"><a href="#">1</a></li>
//      <li><a href="#">2</a></li>
//      <li><a href="#">3</a></li>
//      <li><a href="#">4</a></li>
//      <li><a href="#">5</a></li>
//      <li><a href="#">»</a></li>

// 输出统计的页面
func echoStatisticPages(rowCount int64, curPage, pageSize int) string {
	// math.Ceil()
	pageCount := int(math.Ceil(float64(rowCount) / float64(pageSize)))
	pageshtml := ""
	if curPage <= 1 {
		pageshtml += "<li class=\"am-disabled\"><a href=\"#\">«</a></li>\n"
	} else {
		pageshtml += fmt.Sprintf("<li><a href=\"/admin/statistic?page=%d\">«</a></li>\n", curPage-1)
	}
	// 这里还要优化下，一次只输出多少，后面再以...形式出现
	// 如果总数少于8页，则按正常全部显示
	// 如果当前页大于等于8页了，则第一页固定显示，第二页起以...出现 1...
	for i := 1; i <= pageCount; i++ {
		if i == curPage {
			pageshtml += fmt.Sprintf("<li class=\"am-active\"><a href=\"/admin/statistic?page=%d\">%d</a></li>\n", i, i)
		} else {
			pageshtml += fmt.Sprintf("<li><a href=\"/admin/statistic?page=%d\">%d</a></li>\n", i, i)
		}
	}
	if curPage >= pageCount {
		pageshtml += "<li class=\"am-disabled\"><a href=\"#\">»</a></li>\n"
	} else {
		pageshtml += fmt.Sprintf("<li><a href=\"/admin/statistic?page=%d\">»</a></li>\n", curPage+1)
	}
	return pageshtml
}

func statisticHandler(handler *Handlder) {
	rowCount := queryStatisticCount(handler)
	curPage := handler.QInt("page")
	pageSize := 50
	// 默认没有或者传参数有错误时默认1页
	if curPage <= 0 {
		curPage = 1
	}
	handler.renderTemplate(handler.HTMLFile, ADMINHTML, map[string]interface{}{
		"statistics": queryStatistic(handler, curPage, pageSize),
		"rowCount":   rowCount,
		"pages":      echoStatisticPages(rowCount, curPage, pageSize)})
}
