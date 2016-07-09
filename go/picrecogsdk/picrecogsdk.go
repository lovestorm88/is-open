package picrecogsdk

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const (
	method  = "POST"
	version = "2016-06-29"
)

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

func SignedRequest(uri string) map[string]string {
	return SignedRequest2(uri, PublicKey, PrivateKey, Userid)
}

func SignedRequest2(uri, publicKey, privateKey, userid string) map[string]string {
	params := make(map[string]string)

	params["publicKey"] = publicKey

	params["timestamp"] = time.Now().UTC().Format("2006-01-02T15:04:05Z")

	params["version"] = version

	params["userid"] = userid

	sorted_keys := make([]string, 0)
	for k, _ := range params {
		sorted_keys = append(sorted_keys, k)
	}

	// sort 'string' key in increasing order
	sort.Strings(sorted_keys)

	canonicalized_querys := make([]string, 0, len(params))
	for _, key := range sorted_keys {
		key = strings.Replace(url.QueryEscape(key), "%7E", "~", -1)
		value := strings.Replace(url.QueryEscape(params[key]), "%7E", "~", -1)
		canonicalized_querys = append(canonicalized_querys, fmt.Sprintf("%s=%s", key, value))
	}

	canonicalized_query := strings.Join(canonicalized_querys, "&")

	// create the string to sign
	string_to_sign := method + "\n" + uri + "\n" + canonicalized_query

	// calculate HMAC with SHA256 and base64-encoding
	signature := computeHmac256(string_to_sign, privateKey)

	// encode the signature for the request
	signature = strings.Replace(url.QueryEscape(signature), "%7E", "~", -1)
	//signature = strings.Replace(signature, "%2F", "/", -1)
	params["signature"] = signature

	return params
}

func UploadFileData(url string, params map[string]string, filename string, src io.Reader) (res *http.Response, err error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Add your image file
	fw, err := w.CreateFormFile("file", filename)
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
	req, err := http.NewRequest(method, url, &b)
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
