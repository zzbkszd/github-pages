package xoftp

import (
	"net/http"
	"log"
	"fmt"
	"os"
)

//StartFtpServer desc
func StartFtpServer() {

	err := initDir()

	if(err!=nil){
		fmt.Print(err.Error())
		return
	}

	startHttpServer()


}

func startHttpServer (){

	http.HandleFunc("/interface", func(w http.ResponseWriter,r *http.Request){
		fmt.Fprint(w,"<html><body>")
		fmt.Fprint(w,"<form action='upload' method='post' enctype='multipart/form-data'>")
		fmt.Fprint(w,"<input type='file' name='file'/>")
		fmt.Fprint(w,"<input type='submit' value='提交'/>")
		fmt.Fprint(w,"</form>")
		fmt.Fprint(w,"</body></html>")
	})
	//上传图片
	http.HandleFunc("/upload", ajaxWrapper(upload))


	http.HandleFunc("/editorupload", editorupload)

	//渲染图片
	http.HandleFunc("/render/", render)

	httpDir := http.Dir("upload/")


	//获取图片的方法
	http.Handle("/get/", http.StripPrefix("/get/", http.FileServer(httpDir)))

	fmt.Println(httpDir)

	log.Fatal(http.ListenAndServe(":1179", nil))
}

func initDir() error{
	dir := "upload/"
	finfo, err := os.Stat(dir)
	if err != nil || !finfo.IsDir(){
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			fmt.Printf(err.Error())
		}
		return err
	}
	return nil;
}

