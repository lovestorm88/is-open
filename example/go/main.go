package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/lovestorm88/is-open/go/picrecogsdk"
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
	uri := "/api/porn-recog"

	picrecogsdk.PublicKey = PublicKey
	picrecogsdk.PrivateKey = PrivateKey
	picrecogsdk.Userid = Userid

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("os.Open err")
		return err
	}
	defer file.Close()

	params := picrecogsdk.SignedRequest(uri)

	filename := file.Name()
	res, err := picrecogsdk.UploadFileData(fmt.Sprintf("%s%s", host, uri), params, filename, file)
	if err != nil {
		fmt.Println("UploadFileData err")
		return err
	}

	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("UploadFileData2Gift,ioutil.ReadAll,filename:%s,err:%s", filename, err.Error())
		return err
	}
	defer res.Body.Close()

	/*var data map[string]interface{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		log.Printf("UploadFileData2Gift,Unmarshal,filename:%s,result:%s", filename, result)
		return err
	}*/

	fmt.Println(string(result))

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

	for i := 0; i < *cocurent; i++ {
		wg.Add(1)
		go testPornRecog(&wg, *host, *picPath)
	}

	wg.Wait()
}
