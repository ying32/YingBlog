package blog

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mssola/user_agent"
	//	"github.com/opesun/goquery"
)

// 重定向，方便写函数, M考虑使用mongdb?
type MDB struct {
	*sql.DB
}

// 打开数据库
func openDatabase() *MDB {
	if appSvrConf.DB == "" {
		fmt.Println("数据库还没有设置啦！")
		return nil
	}
	db, err := sql.Open("sqlite3", appSvrConf.DB)
	if err != nil {
		fmt.Println("打开数据库错误，信息:", err.Error())
		return nil
	}
	return &MDB{DB: db}
}

// 关闭数据库
func closeDatabase(db *MDB) {
	if db != nil {
		db.Close()
	}
}

/* **
执行删除操作
-----------------------

param1: 要操作的数据库表
param2: 条件语句，当没有条件时填空
param3: 对应条件传的参数

return 成功返回nil，失败返回error
** */
func (self *MDB) execDelete(tableName, condition string, args ...interface{}) error {
	if self == nil {
		return errors.New("数据库句柄错误。")
	}
	if tableName == "" {
		return errors.New("数据库表名不能为空。")
	}
	tx, err := self.Begin()
	if err != nil {
		return err
	}
	var stmt *sql.Stmt
	if condition == "" || len(args) == 0 {
		stmt, err = tx.Prepare(fmt.Sprintf("DELETE FROM %s", tableName))
	} else {
		stmt, err = tx.Prepare(fmt.Sprintf("DELETE FROM %s WHERE %s", tableName, condition))
	}
	if err != nil {
		return err
	}
	if _, err = stmt.Exec(args...); err != nil {
		return err
	}
	defer stmt.Close()
	if err = tx.Commit(); err != nil {
		// 回滚？？？？用得着么
		tx.Rollback()
		return err
	}
	return nil
}

/* **
执行行查询, 记得添加defer rows.Close()

param1 表名
param2 柱头名
param3 条件
param4 参数
return 数据记录， 错误
** */
func (self *MDB) execQuery(tableName, cols, condition string, args ...interface{}) (rows *sql.Rows, err error) {
	if condition == "" || len(args) == 0 {
		rows, err = self.Query(fmt.Sprintf("SELECT %s FROM %s", cols, tableName))
	} else {
		rows, err = self.Query(fmt.Sprintf("SELECT %s FROM %s WHERE %s", cols, tableName, condition), args...)
	}
	if err != nil {
		return nil, err
	}
	return rows, nil
}

/* **
执行插入命令

param1 表名
param2 柱头，如果不想写入数据填null，否则填 ? 号
param3 参数
return 插入的id，需要有自动编号, 错误
** */
func (self *MDB) execInsert(tableName, cols string, args ...interface{}) (int64, error) {
	tx, err := self.Begin()
	if err != nil {
		return -1, err
	}
	stmt, err := tx.Prepare(fmt.Sprintf("INSERT INTO %s VALUES(%s)", tableName, cols))
	if err != nil {
		return -1, err
	}
	var result sql.Result
	if result, err = stmt.Exec(args...); err != nil {
		return -1, err
	}
	defer stmt.Close()
	if tx.Commit() != nil {
		return -1, err
	}
	insertid, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return insertid, nil
}

/* **
执行更新操作
param1 表名
param2 柱头
parma3 条件
param4 参数
return 影响行数，错误
** */
func (self *MDB) execUpdate(tableName, cols, condition string, args ...interface{}) (int64, error) {
	tx, err := self.Begin()
	if err != nil {
		return 0, err
	}
	if condition == "" || len(args) == 0 {
		return 0, errors.New("条件和参数都不能为空！")
	}
	stmt, err := tx.Prepare(fmt.Sprintf("UPDATE %s SET %s WHERE %s", tableName, cols, condition))
	if err != nil {
		return 0, err
	}
	var result sql.Result
	if result, err = stmt.Exec(args...); err != nil {
		return 0, err
	}
	defer stmt.Close()
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	if affected, err := result.RowsAffected(); err != nil {
		return 0, err
	} else {
		return affected, nil
	}
}

//-----------------以上为公共函数-----------------------

/*
   函数设想：
     利用RTTI反射，定义专有struct，读取专有字段。。。直接填充

*/

// 查询当前用户
func queryCurrentUser(handler *Handlder, user string) (*User, bool) {
	rows, err := handler.DB.Query("select * from users where usr=?", user)
	if err != nil {
		return nil, false
	}
	defer rows.Close()
	usr := &User{}
	if rows.Next() {
		if err = rows.Scan(
			&usr.id,
			&usr.usr,
			&usr.pwd,
			&usr.nickname,
			&usr.isenabled,
			&usr.permission,
			&usr.email,
			&usr.qqnum,
			&usr.address,
			&usr.headimg,
			&usr.signature,
			&usr.lastlogintime,
			&usr.lastloginip); err != nil {
			fmt.Println(err)
			return nil, false
		}
		return usr, true
	}
	return nil, false
}

// 关于Me的信息
func queryAboutMe(handler *Handlder) string {
	rows, err := handler.DB.execQuery("aboutme", "body", "id=1")
	if err != nil {
		return err.Error()
	}
	defer rows.Close()
	if rows.Next() {
		var str string
		if rows.Scan(&str) == nil {
			return str
		}
	}
	return ""
}

// 文章标题信息
type TitlesTop struct {
	Id    int
	Title string
}

// 取最后10条，按时间排
func queryTitlesTop10(handler *Handlder) []TitlesTop {
	//titles := make(map[int]string)
	//rows, err := handler.DB.Query("select id,title from articles order by lastedittime asc limit 10,9")
	var titles []TitlesTop
	rows, err := handler.DB.Query("SELECT id,title FROM articles where ispublic=1 ORDER BY lastedittime DESC LIMIT 0,10")
	if err != nil {
		return titles
	}
	defer rows.Close()
	for rows.Next() {
		//var title string
		//var id int
		var titletop TitlesTop
		rows.Scan(&titletop.Id, &titletop.Title)
		titles = append(titles, titletop)
		//fmt.Println(id, "  ", title)
		//titles[id] = title
	}
	return titles
}

// 文章信息结构
type TArticle struct {
	Id             int
	Title          string
	Author         string
	CreateTime     time.Time
	LastEditTime   time.Time
	HTMLString     string
	MarkdownString string
	Summary        string
	CategoryId     int
	CategoryName   string
	IsPublic       bool
}

func queryArticlesTop5(handler *Handlder) []TArticle {
	var articles []TArticle
	rows, err := handler.DB.Query("select * from articles where ispublic=1 order by lastedittime desc limit 0,5")
	if err != nil {
		return articles
	}
	defer rows.Close()
	for rows.Next() {
		var a TArticle
		var cratetime, lastedittime int64
		rows.Scan(&a.Id, &a.Title, &a.Author, &cratetime, &lastedittime,
			&a.HTMLString, &a.MarkdownString, &a.Summary, &a.CategoryId, &a.IsPublic)
		a.CreateTime = time.Unix(cratetime, 0)
		a.LastEditTime = time.Unix(lastedittime, 0)
		a.HTMLString = base64DeString(a.HTMLString)
		if len(a.HTMLString) > 200 {
			//a.HTMLString = a.HTMLString[:200]
		}
		articles = append(articles, a)
	}
	return articles
}

// 查询文章信息
func queryArticleById(handler *Handlder, id int) *TArticle {

	rows, err := handler.DB.Query("select * from articles where id=?", id)
	if err != nil {
		return nil
	}
	defer rows.Close()
	if rows.Next() {
		a := &TArticle{}
		var cratetime, lastedittime int64
		rows.Scan(&a.Id, &a.Title, &a.Author, &cratetime, &lastedittime,
			&a.HTMLString, &a.MarkdownString, &a.Summary, &a.CategoryId, &a.IsPublic)

		// 要有操作权限的只有当管理员登录后才能查看私有的
		if !a.IsPublic {
			if !checkIsLogin(handler) {
				return nil
			}
		}

		a.CreateTime = time.Unix(cratetime, 0)
		a.LastEditTime = time.Unix(lastedittime, 0)
		a.MarkdownString = base64DeString(a.MarkdownString)
		return a
	}
	return nil
}

// 只查id, 分类, 标题， 作者，最后一次编辑时间
func queryArticlesByPage(handler *Handlder, page, pagesize int) []TArticle {
	var articles []TArticle
	rows, err := handler.DB.Query(
		`select 
		 articles.id,
		 articles.title,
		 articles.author,
		 articles.lastedittime,
		 articles.category,
		 categorys.name, 
		 articles.ispublic
		 from articles 
		 INNER JOIN categorys ON articles.category = categorys.id
		 order by articles.lastedittime 
		 desc limit ?,?`, page*pagesize, pagesize)
	if err != nil {
		fmt.Println(err)
		return articles
	}
	defer rows.Close()
	for rows.Next() {
		var a TArticle
		var lastedittime int64
		rows.Scan(&a.Id, &a.Title, &a.Author, &lastedittime, &a.CategoryId, &a.CategoryName, &a.IsPublic)
		a.LastEditTime = time.Unix(lastedittime, 0)
		articles = append(articles, a)
	}
	return articles
}

// 设置文章是否公开显示
func setArticlePublic(handler *Handlder) bool {
	if aline, err := handler.DB.execUpdate("articles", "ispublic=?", "id=?", handler.PBool("ispublic"), handler.PInt("id")); err == nil && aline > 0 {
		return true
	} else {
		fmt.Println(err, aline)
	}
	return false
}

type ArticleComments struct {
	Id        int
	ArticleId int
	Author    string
	Time      time.Time
	Body      string
	IP        string
}

func queryArticleCommentsById(handler *Handlder, id int) []ArticleComments {
	var comments []ArticleComments
	rows, err := handler.DB.Query("select * from comments where articleid=? order by time desc", id)
	if err != nil {
		fmt.Println(err)
		return comments
	}
	defer rows.Close()
	for rows.Next() {
		var a ArticleComments
		var commenttime int64
		rows.Scan(&a.Id, &a.ArticleId, &a.Author, &commenttime, &a.Body, &a.IP)
		a.Time = time.Unix(commenttime, 0)
		comments = append(comments, a)
	}
	return comments
}

// 插入新的文章返回插入的id，失败为-1
func insertNewArticle(handler *Handlder) (int64, bool) {
	tx, err := handler.DB.Begin()
	if err != nil {
		fmt.Println(err)
		return -1, false
	}
	stmt, err := tx.Prepare("insert into articles values(null,?,?,?,?,?,?,?,?,?)")

	if err != nil {
		fmt.Println(err)
		return -1, false
	}
	var result sql.Result
	//fmt.Println("htmlcode=", handler.PStr("id-html-code"))
	if result, err = stmt.Exec(
		handler.PStr("title"),
		"ying32",
		time.Now().Unix(),
		time.Now().Unix(),
		base64EnString(handler.PStr("id-html-code")),
		base64EnString(handler.PStr("id-markdown-doc")),
		simpleHTMLSummary(handler.PStr("id-html-code"), 400),
		handler.PInt("categroy"),
		1); err != nil {
		fmt.Println(err)
		return -1, false
	}
	defer stmt.Close()
	if tx.Commit() != nil {
		return -1, false
	}
	insertid, err := result.LastInsertId()
	if err != nil {
		return -1, false
	}
	return insertid, true
}

func updateArticle(handler *Handlder) bool {
	tx, err := handler.DB.Begin()
	if err != nil {
		fmt.Println(err)
		return false
	}
	stmt, err := tx.Prepare("update articles set title=?,lastedittime=?,html=?,markdown=?,summary=?,category=? where id=?")

	if err != nil {
		fmt.Println(err)
		return false
	}
	if _, err = stmt.Exec(
		handler.PStr("title"),
		time.Now().Unix(),
		base64EnString(handler.PStr("id-html-code")),
		base64EnString(handler.PStr("id-markdown-doc")),
		simpleHTMLSummary(handler.PStr("id-html-code"), 400),
		handler.PInt("categroy"),
		handler.PUint("id")); err != nil {
		fmt.Println(err)
		return false
	}
	defer stmt.Close()
	if tx.Commit() != nil {
		return false
	}
	return true
}

func deleteArticleById(handler *Handlder, id int) bool {
	tx, err := handler.DB.Begin()
	if err != nil {
		return false
	}
	stmt, err := tx.Prepare("delete from articles where id=?")
	if err != nil {
		return false
	}
	if _, err = stmt.Exec(id); err != nil {
		return false
	}
	defer stmt.Close()
	tx.Commit()
	return true
}

// 分类
type TCategory struct {
	Id   int
	Name string
}

func queryCategroys(handler *Handlder) []TCategory {
	var category []TCategory
	rows, err := handler.DB.Query("select * from categorys")
	if err != nil {
		fmt.Println(err)
		return category
	}
	defer rows.Close()
	for rows.Next() {
		var a TCategory
		rows.Scan(&a.Id, &a.Name)
		category = append(category, a)
	}
	return category
}

func insertStatistic(r *http.Request) {
	// 后台访问暂不计入统计中
	if strings.Contains(r.RequestURI, "/admin") {
		return
	}
	db := openDatabase()
	if db == nil {
		return
	}
	defer closeDatabase(db)
	//db.execInsert("statistic", "null,?,?,?,?,?,?,?,?,?,?,?", )

	tx, err := db.Begin()
	if err != nil {
		return
	}
	stmt, err := tx.Prepare("insert into statistic values(null,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		fmt.Println(err)
		return
	}
	ua := user_agent.New(r.UserAgent())
	/*
	  id, url, referer, refererhost, refererpath, useragent, time, platform, os, browser, browserver ip,
	*/
	refererpath := ""
	refererhost := ""
	referer, err := url.Parse(r.Referer())
	if err == nil {
		refererhost = referer.Host
		refererpath = referer.Path
	}

	browsername, browserversion := ua.Browser()
	if _, err = stmt.Exec(
		r.RequestURI,      // url
		r.Referer(),       // referer
		refererhost,       // refererhost
		refererpath,       // refererpath
		ua.UA(),           // useragent
		time.Now().Unix(), // time
		ua.Platform(),     // platform
		ua.OS(),           // os
		browsername,       // browser
		browserversion,    // browserver
		getIpAddrByStr(r.RemoteAddr)); err != nil {
		fmt.Println(err)
		return
	}
	defer stmt.Close()
	tx.Commit()
	return
}

// 查询的信息，
type TStatistic struct {
	Id             int
	URL            string
	RefererHost    string
	Time           time.Time
	Platform       string
	OS             string
	BrowserName    string
	BrowserVersion string
	IP             string
	Location       string
}

// 返回总数
func queryStatisticCount(handler *Handlder) int64 {
	rows, err := handler.DB.Query("SELECT count(1) FROM statistic")
	if err != nil {
		fmt.Println(err)
		return 0
	}
	defer rows.Close()
	var rowcount int64
	if rows.Next() {
		if err = rows.Scan(&rowcount); err != nil {
			fmt.Println(err)
			return 0
		}
	}
	return rowcount
}

// 查询统计信息
// 		 (SELECT count(statistic.id) FROM statistic) AS rowCount
func queryStatistic(handler *Handlder, curPage, pageSize int) []TStatistic {
	statistics := []TStatistic{}
	rows, err := handler.DB.Query(
		`SELECT 
	     statistic.id,
		 statistic.url,
		 statistic.refererhost,
		 statistic.time,
		 statistic.platform,
		 statistic.os,
		 statistic.browser,
		 statistic.browserver,
		 statistic.ip
	     FROM statistic 
		 order by time
		 desc limit ?,?`, (curPage-1)*pageSize, pageSize)
	if err != nil {
		fmt.Println(err)
		return statistics
	}
	defer rows.Close()
	for rows.Next() {
		var a TStatistic
		var ltime int64
		err = rows.Scan(&a.Id, &a.URL, &a.RefererHost, &ltime, &a.Platform,
			&a.OS, &a.BrowserName, &a.BrowserVersion, &a.IP)
		a.Time = time.Unix(ltime, 0)
		if a.RefererHost != "" {
			a.RefererHost = "http://" + a.RefererHost
		}
		// 这里应该改为插入时将地区写入为好些
		a.Location = QQwryData.GetIPLocationOfString(a.IP)
		statistics = append(statistics, a)
	}
	return statistics
}

// 写个人信息
//func saveUserInfoById(handler *Handlder, usr *User) bool {
//	tx, err := handler.DB.Begin()
//	if err != nil {
//		return false
//	}
//	stmt, err := tx.Prepare("update users set about")
//	if err != nil {
//		return false
//	}
//	defer stmt.Close()
//	tx.Commit()
//}

type TRSSItem struct {
	Id     int
	Title  string
	Author string
	// link省略掉
	Description string
	Comments    string
	PubDate     time.Time
}

func getRSSItems(handler *Handlder) []TRSSItem {
	var items []TRSSItem
	err := handler.DB.execAdvQuery(map[string]interface{}{
		"table":     "articles",
		"cols":      "id,title,author,lastedittime,summary",
		"condition": "ispublic=1",
		"order":     "DESC",
		"ordercol":  "lastedittime",
		"limit":     "0,50"},
		func(rows *sql.Rows) {
			var item TRSSItem
			for rows.Next() {
				var unxitime int64
				if err := rows.Scan(&item.Id, &item.Title, &item.Author, &unxitime, &item.Description); err == nil {
					item.PubDate = time.Unix(unxitime, 0)
					items = append(items, item)
				} else {
					fmt.Println(err)
				}
			}
		})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return items
}
