package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

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
	picPath = flag.String("path", "../resource", "pictures path")
)

func pornRecog(filePath string) error {
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
	res, err := picrecogsdk.UploadFileData(fmt.Sprintf("%s%s", Host, uri), params, filename, file)
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

func testPornRecog(rootPath string) {
	filepath.Walk(rootPath, func(path string, fi os.FileInfo, err error) error {
		if fi == nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}

		err = pornRecog(path)
		if err != nil {
			log.Printf("testPornRecog fail:%s,path:%s\n", err.Error(), path)
		} else {
			log.Printf("testPornRecog success,path:%s\n", path)
		}
		return nil
	})
}

func main() {
	flag.Parse()

	testPornRecog(*picPath)
}
