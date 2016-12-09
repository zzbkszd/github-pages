package xoftp

import (
	"net/http"
	"log"
	"fmt"
	"os"
)

//StartFtpServer desc
func StartFtpServer() {

	initDir()

	startHttpServer()


}


func startHttpServer (){
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/editorupload", editorupload)

	httpDir := http.Dir("./upload/")

	http.Handle("/get/", http.StripPrefix("/get/", http.FileServer(httpDir)))

	fmt.Println(httpDir)

	log.Fatal(http.ListenAndServe(":1179", nil))
}

func initDir() {
	dir := "upload/"
	finfo, err := os.Stat(dir)
	if err != nil || !finfo.IsDir(){
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			fmt.Printf(err.Error())
		}
		return
	}
}

