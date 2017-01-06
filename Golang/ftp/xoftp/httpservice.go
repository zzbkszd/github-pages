package xoftp

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"strconv"
	"bytes"
)

/**
ajax请求接口封装,配置支持跨域
 */
func ajaxWrapper (handler func(http.ResponseWriter,*http.Request) UploadResponse ) func(http.ResponseWriter,*http.Request){

	return func(w http.ResponseWriter,r *http.Request){
		//设置跨域
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Content-Type", "application/json")
		//应对跨域请求的options请求
		if r.Method == "OPTIONS" {
			return
		}

		response := handler(w,r)

		response.Send(w)
	}
}


//上传文件服务
func upload(w http.ResponseWriter, r *http.Request) UploadResponse {

	fmt.Println("file upload active!!")

	//不可使用Get方法
	if r.Method == "GET" {
		return UploadResponse{State: -1, Msg: "请使用post请求上传文件", URL: ""}
	}

	r.ParseMultipartForm(32<<20)


	//获取文件
	file, handler, err := r.FormFile("file")

	if err != nil {
		fmt.Println("get form value file exception : "+err.Error())
		return UploadResponse{State: -1, Msg: err.Error(), URL: ""}
	}

	//提供包支持
	pack := r.FormValue("package")


	var buf bytes.Buffer;
	buf.ReadFrom(file)
	//创建文件结构
	var uploadFile = UploadFile{
		data:buf,
		name:handler.Filename,
		fsroot:"./upload/",
		pack:pack,
	}

	//保存文件
	url,err := uploadFile.save()
	if err != nil {
		return UploadResponse{State: -1, Msg: err.Error(), URL: ""}
	}

	//返回结果
	return UploadResponse{State: 0, Msg: "success", URL: "/get/" + url}

}


/**
为wangeditor专门提供的图片上传接口。
 */
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

	var buf bytes.Buffer;
	buf.ReadFrom(file)
	var uploadFile = UploadFile{
		data:buf,
		name:handler.Filename,
		fsroot:"./upload/"}

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

//对图片进行修改加工
func render(w http.ResponseWriter, r *http.Request) {

	filename := r.URL.Path

	filename = strings.TrimPrefix(filename,"/render/");

	fileHolder := LocalFile{uri:filename}

	err := fileHolder.load()

	if(err!=nil){
		fmt.Fprint(w,"文件不存在")
		return
	}

	scalastr := r.FormValue("scala")
	scalaToStr := r.FormValue("scalaTo")
	cutStr := r.FormValue("cut")

	fmt.Printf(" scalastr=%s \n scalaToStr=%s \n cutStr=%s \n",scalastr,scalaToStr,cutStr)

	//优先使用比例缩放
	if scalastr!="" {

		s,e := strconv.ParseFloat(scalastr,32)
		if(e!=nil){
			return
		}
		s32 := float32(s)
		fileHolder.scala(s32)
	}else if scalaToStr!="" {//无比例缩放数据时按指定缩放

		info := strings.Split(scalaToStr,"*")
		if(len(info)==2){
			x,xerr := strconv.Atoi(info[0])
			y,yerr := strconv.Atoi(info[1])
			if xerr==nil&&yerr==nil {
				fileHolder.scalaAs(x,y)
			}
		}

	}else if cutStr!="" {//裁剪
		info := strings.Split(cutStr,"*")
		if(len(info)==2){
			x,xerr := strconv.Atoi(info[0])
			y,yerr := strconv.Atoi(info[1])
			if xerr==nil&&yerr==nil {
				fmt.Printf("width %d * height %d",x,y)
				fileHolder.cut(x,y)
			}else{
				fmt.Print(xerr)
				fmt.Print(yerr)
			}
		}else{
			fmt.Printf("length of info is %d",len(info))
		}
	}


	fileHolder.httpWrite(w)

	return


	//fmt.Fprintf(w,"request url is ："+filename+"\n")
	//
	//fmt.Fprintf(w,"request values is %s\n",r.FormValue("scala"))
	//fmt.Fprintf(w,"request values is %s\n",r.FormValue("scalaTo"))
	//fmt.Fprintf(w,"request values is %s\n",r.FormValue("cut"))


}