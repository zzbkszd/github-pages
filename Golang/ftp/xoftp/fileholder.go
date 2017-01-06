package xoftp

import(
	"io"
	"os"
	"path/filepath"
	"strings"
	"strconv"
	"crypto/md5"
	"image"
	"ftp/imaging"
	"net/http"
	"github.com/astaxie/beego/logs"
	"fmt"
	"log"
	"bytes"
)

/**
文件保存功能
包括去重，重命名，分包，分组等功能
*/

type UploadFile struct {
	data bytes.Buffer // 文件数据
	name string  //文件名
	md5 string   //md5编码

	urlroot string//url前缀
	fsroot string //文件路径前缀
	pack string //包名

	dupcount int // 查重计数
	issaved int //是否已经保存过：0：未保存，1：已保存

}


/**
主操作，保存文件，并返回文件相对于根路径的访问路径
 */
func (upload *UploadFile) save () (string,error) {

	upload.md5 = calBufferMd5(upload.data)

	//仅当未保存过的时候才保存文件
	if(upload.issaved==0){
		fmt.Printf("save new file to %s\n",upload.fsroot+upload.name)
		//向磁盘存储文件
		f, e := os.OpenFile(upload.fsroot+upload.name, os.O_CREATE|os.O_WRONLY, 0666)
		if(e!=nil){
			log.Printf("error with : %s when open new file %s\n",e.Error(),upload.fsroot+upload.name)
		}
		defer f.Close()
		written,err :=upload.data.WriteTo(f)

		fmt.Printf("save file finish, wirte %d bytes into disk.\n",written)
		//存储完毕
		if err != nil {
			return "",err;
		}
	}

	return upload.urlroot+upload.name,nil
}

//查重，并对重名的不同文件进行重命名
func (upload *UploadFile) checkDuplicate() {
	upload.rename();// 确保文件不重名
}

/**
对文件进行重命名。
 */
func (upload *UploadFile) rename() {
	path := upload.fsroot+upload.name;
	ex := exist(path,upload.md5)
	if ex==1{
		upload.dupcount = upload.dupcount+1 //重复次数+1
		fileext := filepath.Ext(upload.name) //获取扩展名
		filename := strings.TrimSuffix(upload.name,fileext) // 获取文件名
		if(upload.dupcount>1){
			filename = strings.TrimSuffix(upload.name,"_"+strconv.Itoa(upload.dupcount-1)+fileext) // 获取文件名
		}
		filename+="_"+strconv.Itoa(upload.dupcount) //文件名加上重复次数
		upload.name = filename+fileext //重命名
		upload.rename()
	}
	//记录为已保存过
	if ex==2 {
		upload.issaved=1
	}
	return
}



// 检查文件是否存在
// 根据MD5判断重名文件是否重复，若重复则删除原文件
// 如果由 filename 指定的文件或目录存在则返回 1，否则返回0，若文件已存在且MD5相等，返回2
func exist(filename string,md5 string) int {
	_, err := os.Stat(filename)
	//文件不存在
	if err != nil {
		return 0
	}
	file, ferr := os.Open(filename)
	//读取文件失败=不存在
	if ferr!=nil {
		return 0
	}
	//若MD5相等，保存原文件,算作已存在
	if calMd5(file) == md5{
		return 2;
	}

	//文件存在
	return 1;
}

//计算MD5值
func calMd5(input io.Reader) string {
	md5h := md5.New()
	io.Copy(md5h, input)
	return string(md5h.Sum([]byte(""))) //md5
}
//计算MD5值
func calBufferMd5(input bytes.Buffer) string {
	md5h := md5.New()
	input.WriteTo(md5h)
	return string(md5h.Sum([]byte(""))) //md5
}

/**
本地图片文件
 */
type LocalFile struct {
	uri string
	data image.Image//图片数据
	out string //输出文件路径
	scalaTo float32 //缩放比例
	scalaAsStr string//缩放尺寸字符串
	scalaWidth int
	scalaHeight int
	cutStr string //裁剪属性字符串
	cutStartX int
	cutStartY int
	cutWidth  int
	cutHeight int

}

//将本地文件加载到内存,若读取文件出错则返回错误
func (file *LocalFile) load () error{
	path := "upload/"+file.uri
	logs.Info(path)
	f,err := os.Open(path)
	if(err!=nil){
		return err
	}
	defer f.Close()

	img,ie := imaging.Decode(f)

	if(ie!=nil){
		return ie
	}

	file.data = img

	return nil

}

//按比例缩放
func (file *LocalFile) scala(s float32){

	bound := file.data.Bounds()

	dx := float32(bound.Dx())
	dy := float32(bound.Dy())

	fmt.Printf("before scala image size is %f*%f\n",dx,dy)

	file.scalaWidth = int(dx*s)
	file.scalaHeight = int(dy*s)

	fmt.Printf("scala image as %f, to width %d * height %d\n",s,file.scalaWidth,file.scalaHeight)

	dst := imaging.Resize(file.data, file.scalaWidth, file.scalaHeight, imaging.Lanczos)

	file.data = dst

}

func (file *LocalFile) scalaAs(x,y int){
	fmt.Printf("scala image to %d * %d\n",x,y)
	file.scalaWidth = x
	file.scalaHeight = y
	dst := imaging.Fit(file.data, file.scalaWidth, file.scalaHeight, imaging.Lanczos)
	file.data = dst
}

func (file *LocalFile) cut(x,y int){
	fmt.Printf("cut image to %d * %d\n",x,y)
	file.cutHeight = x
	file.cutWidth = y
	dst := imaging.Fill(file.data,x,y,imaging.Center,imaging.Lanczos)
	file.data = dst
}

func (file *LocalFile) httpWrite (w http.ResponseWriter){
	w.Header().Set("contentType","image/png")
	imaging.Encode(w,file.data,imaging.PNG)
}

