package main

import (
	// "fmt"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

var count int

// CheckAndCreateDir checks if a directory exists and creates it if it doesn't
func CheckAndCreateDir(path string) error {
	// Check if the directory exists
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// Directory does not exist, create it
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
		fmt.Printf("Directory created: %s\n", path)
	} else if err != nil {
		// Other errors
		return fmt.Errorf("error checking directory: %v", err)
	} else {
		fmt.Printf("Directory already exists: %s\n", path)
	}
	return nil
}
func htmltext(url string) string {
	client := http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	// fmt.Println(string(body))
	return string(body)
}
func urlList(s string) ([]string, string) {
	html := htmltext(s)
	tmp := strings.ReplaceAll(html, "\n", "")
	// fmt.Println(tmp)
	ex := regexp.MustCompile("<img src=\"/file/(.*?)\">")
	flodernameEx := regexp.MustCompile("<title>(.*?)</title>")
	r := ex.FindAllStringSubmatch(tmp, -1)
	ret := []string{}

	for _, i := range r {
		ret = append(ret, i[1])
	}
	count = len(ret)
	flodername := flodernameEx.FindAllStringSubmatch(tmp, -1)
	CheckAndCreateDir(`your path` + flodername[0][1])
	// fmt.Println(len(flodername))
	return ret, flodername[0][1]
}
func downloadImage(imgaename string, flodername string) {

	req, err := http.Get(`https://telegra.ph/file/` + imgaename)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()
	content, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	entries, err := os.ReadDir(`your path` + flodername + "/")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 计算文件数量
	nFiles := len(entries)
	for i := 0; i <= nFiles; i++ {
		fmt.Printf(".")
	}

	fmt.Printf("下载进度%.2f%%\n", (float64(nFiles)/float64(count))*100)
	cmd := exec.Command("clear")
	cmd.Stdout = exec.Command("clear").Stdout
	cmd.Run()
	err = os.WriteFile(`your path`+flodername+"/"+imgaename, content, 0666)
	if err != nil {
		panic(err)
	}
}
func printSliceConcurrently(imagelist []string, flodername string, wg *sync.WaitGroup) {
	for _, num := range imagelist {
		wg.Add(1) // 增加等待组计数
		go func(image string) {
			defer wg.Done() // goroutine结束时减少计数
			downloadImage(image, flodername)
		}(num)
	}
	wg.Wait() // 等待所有goroutine完成
}
func main() {

	var s string
	fmt.Println("输入链接：")
	fmt.Scanln(&s)
	r, f := urlList(s)
	var wg sync.WaitGroup

	printSliceConcurrently(r, f, &wg)
	for i := 0; i <= count; i++ {
		fmt.Printf(".")
	}
	fmt.Println("Success")
}
