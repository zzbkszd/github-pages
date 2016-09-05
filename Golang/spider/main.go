package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
)

var (
	//图片正则表达式
	imageItemExp = regexp.MustCompile(`src="//i\.4cdn\.org/s/[0123456789]+s\.jpg"`)
	//帖子路径正则表达式
	threadItemExp = regexp.MustCompile(`"thread/[0123456789]+"`)
)

//ThreadItem 帖子数据
type ThreadItem struct {
	url     string   //帖子路径
	content string   //帖子内容
	imgs    []string //帖子图片
}

//获取网页内容
func (t *ThreadItem) getContent() *ThreadItem {
	content, err := httpGet(t.url)
	if err != 200 {
		t.content = ""
		return t
	}
	t.content = string(content)
	return t
}

//从网页内容中抓取图片链接
func (t *ThreadItem) getImage() *ThreadItem {
	imgs := imageItemExp.FindAllStringSubmatch(t.content, 10000)
	var l = make([]string, 0)
	for _, v := range imgs {
		l = append(l, v[0])
	}
	t.imgs = l
	return t
}

//下载所有抓取的图片链接
func (t *ThreadItem) download() {
	last := strings.LastIndex(t.url, "/")
	dir := "download/" + string(t.url[last:len(t.url)])
	fmt.Println("create dir:", dir)
	err := os.Mkdir(dir, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, img := range t.imgs {
		pos := strings.LastIndex(img, "/")
		filename := string(img[pos : len(img)-1])
		file, err := os.Create(dir + "/" + filename)
		defer file.Close()
		if err != nil {
			fmt.Println("error for create file")
			fmt.Println(err)
			continue
		}
		data, error := downloadImg("http:" + string(img[5:len(img)-1]))
		if error != 200 {
			fmt.Println("error for download image:", error)
			continue
		}
		file.Write(data)
	}
}

/*
找到帖子链接
*/
func findThreads(url string) []ThreadItem {
	var threads = make([]ThreadItem, 0)
	content, err := httpGet(url)
	if err != 200 {
		return threads
	}
	tds := threadItemExp.FindAllStringSubmatch(content, 10000)
	var tdStr = make([]string, 0)
	//去掉引号，并放到一维数组中
	for _, t := range tds {
		var n = strings.Replace(t[0], "\"", "", -1)
		tdStr = append(tdStr, n)
	}
	//去重准备
	sort.Strings(tdStr)
	tdStr = unequal(tdStr)
	//组装帖子结构体
	for _, t := range tdStr {
		threads = append(threads, ThreadItem{url: "http://boards.4chan.org/s/" + t})
	}
	return threads
}

func downloadImg(url string) (content []byte, statusCode int) {
	url = strings.Replace(url, "s.", ".", -1)
	fmt.Println("download img from url:", url)
	resp, err1 := http.Get(url)
	if err1 != nil {
		statusCode = -100
		return
	}
	if resp.StatusCode == 404 {
		url = strings.Replace(url, ".jpg", ".png", -1)
		resp, err1 = http.Get(url)
		if err1 != nil {
			statusCode = -100
			return
		}
	}
	defer resp.Body.Close()
	content, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		statusCode = -200
		return
	}
	statusCode = resp.StatusCode
	return
}

/*
http获取方法
*/
func httpGet(url string) (content string, statusCode int) {
	resp, err1 := http.Get(url)
	if err1 != nil {
		statusCode = -100
		return
	}
	defer resp.Body.Close()
	data, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		statusCode = -200
		return
	}
	statusCode = resp.StatusCode
	content = string(data)
	return
}

/*
去重
*/
func unequal(a []string) (ret []string) {
	aLen := len(a)
	for i := 0; i < aLen; i++ {
		if i > 0 && a[i-1] == a[i] {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}

/*
爬虫入口
*/
func work(url string) {
	fmt.Println("get list with url :", url)
	var threads = findThreads(url)
	fmt.Println(threads)
	for _, v := range threads {
		(&v).getContent().getImage().download()
		// fmt.Println(v.imgs)
	}
}

func main() {
	// work("http://boards.4chan.org/s/")
	pages := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10"}
	for _, index := range pages {
		work("http://boards.4chan.org/s/" + index + "/")
	}
}
