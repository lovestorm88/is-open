package sdk

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const (
	METHOD  = "POST"
	VERSION = "1"
)

//图片识别类型
const (
	PIC_RECOG_PORN = "/api/porn-recog"
)

//公共结果部分
type CommonRsp struct {
	ErrCode int32  `json:"errCode"`
	Msg     string `json:"msg"`
}

//图片检测
type PicRecogRsp struct {
	CommonRsp
	Name       string  `json:"name"`
	Label      int32   `json:"label"`
	Confidence float64 `json:"confidence"`
}

//批量检测
type BatchPicRecogRsp struct {
	CommonRsp
	Data []PicRecogRsp `json:"data"`
}

var (
	PublicKey  string
	PrivateKey string
	Userid     string
)

func computeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func SignedRequest() map[string]string {
	return signedRequest(PublicKey, PrivateKey, Userid)
}

func signedRequest(publicKey, privateKey, userid string) map[string]string {
	params := make(map[string]string)

	params["publicKey"] = publicKey

	params["timestamp"] = fmt.Sprintf("%d", time.Now().Unix())

	params["version"] = VERSION

	params["userid"] = userid

	sorted_keys := make([]string, 0)
	for k, _ := range params {
		sorted_keys = append(sorted_keys, k)
	}

	// sort 'string' key in increasing order
	sort.Strings(sorted_keys)

	canonicalized_querys := make([]string, 0, len(params))
	for _, key := range sorted_keys {
		key = url.QueryEscape(key)
		value := url.QueryEscape(params[key])
		canonicalized_querys = append(canonicalized_querys, fmt.Sprintf("%s=%s", key, value))
	}

	// create the string to sign
	string_to_sign := strings.Join(canonicalized_querys, "&")

	// calculate HMAC with SHA256 and base64-encoding
	signature := computeHmac256(string_to_sign, privateKey)
	//fmt.Printf("string_to_sign:%s,privateKey:%s,signature:%s,", string_to_sign, privateKey, signature)

	// encode the signature for the request
	signature = url.QueryEscape(signature)
	//fmt.Printf("signature_escape:%s\n", signature)
	params["signature"] = signature

	return params
}

func UploadFileData(url string, params map[string]string, filename string, src io.Reader) (res *http.Response, err error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Add your image file
	fw, err := w.CreateFormFile("image", filename)
	if err != nil {
		return
	}
	if _, err = io.Copy(fw, src); err != nil {
		return
	}

	// Add the other fields
	for k, v := range params {
		if fw, err = w.CreateFormField(k); err != nil {
			return
		}

		if _, err = fw.Write([]byte(v)); err != nil {
			return
		}
	}

	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest(METHOD, url, &b)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	client := &http.Client{}
	res, err = client.Do(req)
	if err != nil {
		return
	}

	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
	}

	return
}

func PicRecog(host string, picRecogType string, filename string, file io.Reader) (*BatchPicRecogRsp, error) {
	params := signedRequest(PublicKey, PrivateKey, Userid)
	res, err := UploadFileData(fmt.Sprintf("%s%s", host, picRecogType), params, filename, file)
	if err != nil {
		fmt.Printf("PicRecog,err1:%s\n", err.Error())
		return nil, err
	}

	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("PicRecog,err2:%s\n", err.Error())
		return nil, err
	}
	defer res.Body.Close()

	fmt.Println(string(result))
	var brsp BatchPicRecogRsp
	err = json.Unmarshal(result, &brsp)
	if err != nil {
		fmt.Printf("PicRecog,err3:%s\n", err.Error())
		return nil, err
	}

	return &brsp, nil
}
