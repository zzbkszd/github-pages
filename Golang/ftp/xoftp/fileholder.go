package xoftp

import(
	"io"
	"os"
	"path/filepath"
	"strings"
	"strconv"
	"crypto/md5"
)

/**
文件保存功能
包括去重，重命名，分包，分组等功能
*/

type UploadFile struct {
	data io.Reader // 文件数据
	name string  //文件名
	md5 string   //md5编码
	urlroot string//url前缀
	fsroot string //文件路径前缀
	dupcount int // 查重计数
	issaved int //是否已经保存过：0：未保存，1：已保存
}

/**
主操作，保存文件，并返回文件相对于根路径的访问路径
 */
func (upload *UploadFile) save () (string,error) {

	upload.md5 = calMd5(upload.data)
	upload.checkDuplicate();

	//仅当未保存过的时候才保存文件
	if(upload.issaved==0){
		//向磁盘存储文件
		f, _ := os.OpenFile(upload.fsroot+upload.name, os.O_CREATE|os.O_WRONLY, 0660)
		_,error := io.Copy(f, upload.data)
		//存储完毕
		if error != nil {
			return "",error;
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
	if ex==-1{
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

// 检查文件或目录是否存在
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
	//若MD5相等，保存原文件
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
