package blog

// 此单元是参照gopher来写，所以文件名和函数名基本一样

import (
	"fmt"
	"net/http"
	"os"
	//"strings"
	"time"

	"github.com/gorilla/mux"
)

func handlerFunc(route Route) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//if strings.Contains(r.UserAgent(), "Python-urllib") {
		//	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		//	return
		//}

		handler := NewHandler(w, r, route.HTMLFile)
		// 关
		defer closeDatabase(handler.DB)

		// 等下这里要做异步处理
		go insertStatistic(r)

		url := time.Now().Format(time.RFC850) + " " + r.Method + " " + r.URL.Path
		if r.URL.RawQuery != "" {
			url += "?" + r.URL.RawQuery
		}
		fmt.Println(url)
		// 任何人都可以访问的
		if route.Permission&Everyone == Everyone {
			route.HandleFunc(handler)
			return
		}
		// 管理员页面访问
		// 注释掉，方便调式
		if route.Permission&Administrator == Administrator {
			if !checkIsLogin(handler) {
				redirectLogin(handler)
				return
			}
		}
		route.HandleFunc(handler)
	}
}

// 启动服务
func StartBlogServer() {

	// 生成上传文件夹
	if !isExists(appSvrConf.UploadPath) {
		//os.Mkdir(WEB_DEFAULT_UPLOADS_PATH, 0)
		os.MkdirAll(appSvrConf.UploadPath+"images/", 0777)
	}

	r := mux.NewRouter()
	for _, route := range routes {
		r.HandleFunc(route.URL, handlerFunc(route))
	}
	// 自定义404
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	// 下面两种方式各有不同
	// 以下方式通过函数回调来加载文件，自己管理
	//r.PathPrefix("/assets/").HandlerFunc(fileHandler)
	// 以下方式等同于http.Handle("/", http.FileServer(http.Dir("./AmazeUI/assets")))的效果直接静态了整个目录交由go去管理了
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	// 开源代码映射
	//http.Handle("/sources/", http.StripPrefix("/sources/", http.FileServer(http.Dir("./blog"))))
	//http.Handle("/sources/template/", http.StripPrefix("/sources/template/", http.FileServer(http.Dir("./template"))))
	// 映射上传目录
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir(appSvrConf.UploadPath))))

	http.Handle("/", r)

	fmt.Println("服务已启动。")
	if appSvrConf.Port == 80 {
		fmt.Printf("请在浏览器中输入: http://%s/来访问主页\n", getLocalIp())
	} else {
		fmt.Printf("请在浏览器中输入: http://%s:%d/来访问主页\n", getLocalIp(), appSvrConf.Port)
	}
	err := http.ListenAndServe(fmt.Sprintf(":%d", appSvrConf.Port), nil)
	if err != nil {
		panic(err)
	}
}
