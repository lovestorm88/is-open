package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	sdk "github.com/lovestorm88/is-open/sdk/go"
)

//以下参数需要根据实际情况修改
const (
	PublicKey  = "publicKey_test"
	PrivateKey = "privateKey_test"
	Userid     = "userid_test"
	Host       = "http://localhost:8087"
)

var (
	host     = flag.String("host", "http://localhost:8087", "host url")
	picPath  = flag.String("path", "../resource", "pictures path")
	cocurent = flag.Int("cocurrent", 1, "cocurrent number")
)

func pornRecog(host, filePath string) error {
	sdk.PublicKey = PublicKey
	sdk.PrivateKey = PrivateKey
	sdk.Userid = Userid

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("os.Open err")
		return err
	}
	defer file.Close()

	filename := file.Name()
	brsp, err := sdk.PicRecog(host, sdk.PIC_RECOG_PORN, filename, file)
	if err != nil {
		fmt.Println("UploadFileData err")
		return err
	}

	fmt.Println(brsp)

	return err
}

func testPornRecog(wg *sync.WaitGroup, host, rootPath string) {
	defer wg.Done()

	filepath.Walk(rootPath, func(path string, fi os.FileInfo, err error) error {
		if fi == nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}

		st := time.Now().UnixNano()
		err = pornRecog(host, path)
		et := time.Now().UnixNano()
		if err != nil {
			log.Printf("testPornRecog fail:%s,use:%d,path:%s\n", err.Error(), et-st, path)
		} else {
			log.Printf("testPornRecog success,use:%d,path:%s\n", et-st, path)
		}
		return nil
	})
}

func main() {
	flag.Parse()
	var wg sync.WaitGroup

	st := time.Now().Unix()
	for i := 0; i < *cocurent; i++ {
		wg.Add(1)
		go testPornRecog(&wg, *host, *picPath)
	}

	wg.Wait()
	et := time.Now().Unix()
	log.Printf("total use:%d\n", et-st)
}
