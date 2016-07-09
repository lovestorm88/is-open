package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/lovestorm88/is-open/go/picrecogsdk"
)

func testPornRecog1(path string) error {
	var (
		host = "http://localhost:8087"
		uri  = "/api/porn-recog"
	)

	picrecogsdk.PublicKey = "publicKey_test"
	picrecogsdk.PrivateKey = "privateKey_test"
	picrecogsdk.Userid = "userid_test"

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	params := picrecogsdk.SignedRequest(uri)

	filename := file.Name()
	res, err := picrecogsdk.UploadFileData(fmt.Sprintf("%s%s", host, uri), params, filename, file)
	if err != nil {
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

func main() {
	path := "../resource/1.jpg"
	err := testPornRecog1(path)
	if err != nil {
		log.Printf("testPornRecog1 fail:%s,path:%s\n", err.Error(), path)
	} else {
		log.Printf("testPornRecog1 success,path:%s\n", path)
	}

}
