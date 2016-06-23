// YingBlog project main.go
package main

import (
	"YingBlog/blog"
	"runtime"
)

func main() {
	// 多核设置
	runtime.GOMAXPROCS(runtime.NumCPU())
	blog.StartBlogServer()
}
