package main

import (
	"flag"
	"fmt"
	"io"
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
	batch    = flag.Int("batch", 1, "batch number")
)

func pornRecog(host, filePath string) int32 {
	sdk.PublicKey = PublicKey
	sdk.PrivateKey = PrivateKey
	sdk.Userid = Userid

	file, err := os.Open(filePath)
	if err != nil {
		return -1
	}
	defer file.Close()

	filenames := make([]string, 0, *batch)
	files := make([]io.Reader, 0, *batch)
	filename := file.Name()
	for i := 0; i < *batch; i++ {
		filenames = append(filenames, filename)
		files = append(files, file)
	}

	brsp, err := sdk.BatchPicRecog(host, sdk.PIC_RECOG_PORN, filenames, files)
	if err != nil {
		return -2
	}

	return brsp.ErrCode
}

func testPornRecog(index int, wg *sync.WaitGroup, host, rootPath string) {
	defer wg.Done()

	errCodes := make(map[int32]int)

	filepath.Walk(rootPath, func(path string, fi os.FileInfo, err error) error {
		if fi == nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}

		errCode := pornRecog(host, path)
		errCodes[errCode] += 1

		return nil
	})

	log.Printf("index:%d,result:%v", index, errCodes)

}

func main() {
	flag.Parse()
	var wg sync.WaitGroup

	st := time.Now().Unix()
	for i := 0; i < *cocurent; i++ {
		wg.Add(1)
		go testPornRecog(i, &wg, *host, *picPath)
	}

	wg.Wait()
	et := time.Now().Unix()
	log.Printf("total use:%d\n", et-st)
}
