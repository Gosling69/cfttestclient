package requests

import (
	"bytes"
	"cfttestclient/workwithfiles"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

var ServerHost = ""
var ServerPort = ""

func Request(url, method string, attachment ...string) (int, []byte) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	if len(attachment) > 0 {
		fileName := attachment[0]
		file, errFile1 := os.Open(fileName)
		if errFile1 != nil {
			log.Println(errFile1)
		}
		defer file.Close()
		part1, errFile1 := writer.CreateFormFile("file", fileName)
		if errFile1 != nil {
			log.Println(errFile1)
		}
		_, errFile1 = io.Copy(part1, file)
		if errFile1 != nil {
			log.Println(errFile1)
		}
		err := writer.Close()
		if err != nil {
			log.Println(err)
		}
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return 0, []byte{}
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err != nil {
		fmt.Println(err)
		return 0, []byte{}
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return 0, []byte{}
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return 0, []byte{}
	}
	return res.StatusCode, body
}

func GetFiles() (map[string]interface{}, bool) {
	result := map[string]interface{}{}
	url := fmt.Sprintf("http://%v:%v/getFileList", ServerHost, ServerPort)
	method := "GET"
	code, body := Request(url, method)
	if code == 0 {
		return result, false
	}
	json.Unmarshal(body, &result)
	return result, true
}
func GetFile(name string) bool {
	url := fmt.Sprintf("http://%v:%v/getFile/%v", ServerHost, ServerPort, name)
	method := "GET"
	code, body := Request(url, method)
	if code == http.StatusOK {
		return workwithfiles.SaveFile(name, body)
	}
	return code == http.StatusOK
}
func UpdateFile(name string) bool {
	url := fmt.Sprintf("http://%v:%v/updateFile", ServerHost, ServerPort)
	method := "POST"
	code, body := Request(url, method, name)
	if code == http.StatusOK && len(body) > 0 {
		fmt.Println(string(body))
	}
	return code == http.StatusOK
}
func PutFile(name string) bool {
	url := fmt.Sprintf("http://%v:%v/putFile", ServerHost, ServerPort)
	method := "PUT"
	code, _ := Request(url, method, name)
	return code == http.StatusOK
}

func DeleteFile(name string) bool {
	url := fmt.Sprintf("http://%v:%v/deleteFile/%v", ServerHost, ServerPort, name)
	method := "DELETE"
	code, _ := Request(url, method)
	return code == http.StatusOK
}
