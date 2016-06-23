package blog

import (
	"fmt"
	"io"
	"os"
	"time"
)

type imageUpload struct {
	Success int    `json:"success"`
	Message string `json:"message"`
	URL     string `json:"url"`
}

func uploadHandler(handler *Handlder) {
	if handler.Request.Method == "POST" {
		file, fileHeader, err := handler.Request.FormFile("editormd-image-file")
		if err != nil {
			data := &imageUpload{0, "上传图片失败,  错误：" + err.Error(), ""}
			handler.renderJson(data)
			return
		}
		fmt.Println("上传的文件为：", fileHeader.Filename)
		if !isCanUploadFile(fileHeader.Filename) {
			data := &imageUpload{0, "不允许上传此类型文件", ""}
			handler.renderJson(data)
			return
		}
		defer file.Close()
		// 产生的新文件
		newFileName := fmt.Sprintf("images/%d%s", time.Now().UnixNano(), getFileExt(fileHeader.Filename))
		// 文件绝对路径
		newFile := fmt.Sprintf("%s%s", appSvrConf.UploadPath, newFileName)
		// 相对路么文件名
		newRelativelyFileName := fmt.Sprintf("/resources/%s", newFileName)

		fmt.Println("new File:", newFile)
		f, err := os.Create(newFile)
		if err != nil {
			data := &imageUpload{0, "上传图片失败,  错误：" + err.Error(), ""}
			handler.renderJson(data)
			return
		}
		defer f.Close()
		_, err = io.Copy(f, file)
		if err != nil {
			data := &imageUpload{0, "上传图片失败,  错误：" + err.Error(), ""}
			handler.renderJson(data)
			return
		}
		data := &imageUpload{1, "上传图片成功！", newRelativelyFileName}
		handler.renderJson(data)
	}
}
