package xoftp

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//UploadResponse 返回消息
type UploadResponse struct {
	State int    `json:"state"`
	URL   string `json:"data"`
	Msg   string `json:"msg"`
}

//SendJsonp 发送jsonp消息
func (ur *UploadResponse) SendJsonp(w http.ResponseWriter) {
	fmt.Println("send message msg:", ur.Msg, " state : ", ur.State)
	doc, err := json.Marshal(ur)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Fprintf(w, string(doc))
}

//Send 发送消息
func (ur *UploadResponse) Send(w http.ResponseWriter) {
	fmt.Println("send message msg:", ur.Msg, " state : ", ur.State)
	doc, err := json.Marshal(ur)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Fprintf(w, string(doc))
}
