package xoftp

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

//StartFtpServer desc
func StartFtpServer() {

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string("hello"))
	})
	http.HandleFunc("/upload", upload)

	httpDir := http.Dir("./upload/")

	http.Handle("/get/", http.StripPrefix("/get/", http.FileServer(httpDir)))

	fmt.Println(httpDir)

	log.Fatal(http.ListenAndServe(":1179", nil))

}

func upload(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "GET" {
		var response = UploadResponse{State: -1, Msg: "请使用post请求上传文件", URL: ""}
		response.Send(w)
		return
	}
	if r.Method == "OPTIONS" {
		return
	}

	// Stop here if its Preflighted OPTIONS request

	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		var response = UploadResponse{State: -1, Msg: err.Error(), URL: ""}
		response.Send(w)
		return
	}

	//文件扩展名
	fileext := filepath.Ext(handler.Filename)
	//用时间戳来保证文件名不重复
	//TODO 可以在此做一些去重的工作
	filename := strconv.FormatInt(time.Now().Unix(), 10) + fileext

	f, _ := os.OpenFile("./upload/"+filename, os.O_CREATE|os.O_WRONLY, 0660)
	_, err = io.Copy(f, file)
	if err != nil {
		var response = UploadResponse{State: -2, Msg: err.Error(), URL: ""}
		response.Send(w)
		return
	}
	filedir, _ := filepath.Abs("./upload/" + filename)
	fmt.Println(filename + "上传完成,服务器地址:" + filedir)
	var response = UploadResponse{State: 0, Msg: "success", URL: "/get/" + filename}
	response.Send(w)
	return

}

func initDir() {
	dir := "upload/"
	err := os.Mkdir(dir, os.ModePerm)
	if err != nil {
		fmt.Printf(err.Error())
	}
}
