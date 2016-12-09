package xoftp

import (
	"fmt"
	"net/http"
	"path/filepath"
)

//上传文件服务
func upload(w http.ResponseWriter, r *http.Request) {

	//设置跨域
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Content-Type", "application/json")

	fmt.Println("file upload active!!")

	//不可使用Get方法
	if r.Method == "GET" {
		var response = UploadResponse{State: -1, Msg: "请使用post请求上传文件", URL: ""}
		response.Send(w)
		return
	}
	//应对跨域请求的options请求
	if r.Method == "OPTIONS" {
		return
	}

	//获取文件
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		var response = UploadResponse{State: -1, Msg: err.Error(), URL: ""}
		response.Send(w)
		return
	}

	//创建文件结构
	var uploadFile = UploadFile{data:file,name:handler.Filename,fsroot:"./upload/"}


	//保存文件
	url,err := uploadFile.save()
	if err != nil {
		var response = UploadResponse{State: -1, Msg: err.Error(), URL: ""}
		response.Send(w)
		return
	}

	//返回结果
	var response = UploadResponse{State: 0, Msg: "success", URL: "/get/" + url}
	response.Send(w)
	return

}

func editorupload(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if r.Method == "GET" {
		//var response = UploadResponse{State: -1, Msg: "请使用post请求上传文件", URL: ""}
		fmt.Fprintf(w,"error|请求错误")
		return
	}
	if r.Method == "OPTIONS" {
		return
	}

	// Stop here if its Preflighted OPTIONS request

	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("wangEditorH5File")
	if err != nil {
		fmt.Println("read file : "+err.Error());
		fmt.Fprintf(w,"error|"+err.Error())
		return
	}

	var uploadFile = UploadFile{data:file,name:handler.Filename}

	_,e := uploadFile.save()
	if e != nil {
		var response = UploadResponse{State: -1, Msg: e.Error(), URL: ""}
		response.Send(w)
		return
	}

	filedir, _ := filepath.Abs("./upload/" +  uploadFile.name)
	fmt.Println( uploadFile.name + "上传完成,服务器地址:" + filedir)

	fmt.Fprintf(w,"http://123.206.43.110:1179/get/" + uploadFile.name);

	return

}
