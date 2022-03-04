package main

import (
	"bufio"
	"cfttestclient/requests"
	"cfttestclient/workwithfiles"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Config struct {
	ServerHost string `json:"serverHost"`
	ServerPort string `json:"serverPort"`
	SaveFolder string `json:"saveFolder"`
}

func ShowCommands() string {
	return `
get <filename> - get file and write to folder specified in config
put <filename> - upload file to server
update <filename> - update file on server
delete <filename> - delete file from server
help - show this message again
exit - close client
`
}

func InputHostPort(reader *bufio.Reader) {
	fmt.Println("Input server host ->")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	requests.ServerHost = text
	fmt.Println("Input server port ->")
	text, _ = reader.ReadString('\n')
	text = strings.TrimSpace(text)
	requests.ServerPort = text
}

func HandleInput(input string) {
	var result interface{}
	command := strings.Split(strings.TrimSpace(input), " ")
	if len(command) >= 2 {
		filename := strings.Join(command[1:], " ")
		filename = strings.TrimSpace(filename)
		switch command[0] {
		case "get":
			result = requests.GetFile(filename)
		case "put":
			if exists, _ := workwithfiles.FileExists(filename); exists {
				result = requests.PutFile(filename)
			} else {
				result = "No such file"
			}
		case "update":
			if exists, _ := workwithfiles.FileExists(filename); exists {
				result = requests.UpdateFile(filename)
			} else {
				result = "No such file"
			}
		case "delete":
			result = requests.DeleteFile(filename)
		default:
			result = "unknown command"
		}
	} else {
		switch command[0] {
		case "help":
			result = ShowCommands()
		case "exit":
			os.Exit(0)
		default:
			result = "Incorrect input"
		}
	}
	fmt.Println(result)
}

func ReadConfig() error {
	jsonFile, err := os.Open("conf.json")
	if err != nil {
		// fmt.Println(err)
		return err
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var config Config
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(config.SaveFolder, 0777); err == nil {
		workwithfiles.PathBase = config.SaveFolder
	}
	requests.ServerHost = config.ServerHost
	requests.ServerPort = config.ServerPort
	return nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	confErr := ReadConfig()
	if confErr != nil {
		fmt.Println("Couldn't read config")
		InputHostPort(reader)
	}
	fileList, serverOk := requests.GetFiles()
	for !serverOk {
		fmt.Println("Server unreachable")
		InputHostPort(reader)
		fileList, serverOk = requests.GetFiles()
	}
	fmt.Println(ShowCommands())
	fmt.Println("Files on server:")
	for file, hash := range fileList {
		fmt.Printf("%v - %v\n", file, hash)
	}
	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)
		HandleInput(text)
	}
}
