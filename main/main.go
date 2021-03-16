package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var telegramKey string = ""
var chatID int64 = 0

func telegramBaseURL() string {
	return "https://api.telegram.org/bot" + telegramKey
}

func runBot() {
	updatesChan := make(chan updatesModel)
	go getUpdates(updatesChan, 0, 0)
	for true {
		updates := <-updatesChan
		processResponse(updates)
	}
}

func getUpdates(updates chan updatesModel, offset int64, sleepDuration time.Duration) {
	time.Sleep(sleepDuration)
	url := telegramBaseURL() + "/getUpdates"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		getUpdates(updates, offset, sleepDuration+2)
		return
	}

	q := req.URL.Query()
	q.Add("timeout", "100")
	if offset > 0 {
		q.Add("offset", strconv.FormatInt(offset, 10))
	}
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		getUpdates(updates, offset, sleepDuration+2)
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		getUpdates(updates, offset, sleepDuration+2)
		return
	}

	value := updatesModel{}
	err = json.Unmarshal(body, &value)
	if err != nil {
		getUpdates(updates, offset, sleepDuration+2)
		return
	}

	updates <- value
	if len(value.Result) > 0 {
		offset = value.Result[len(value.Result)-1].UpdateId + 1
	}
	getUpdates(updates, offset, 0)
}

func processResponse(model updatesModel) {
	if len(model.Result) == 0 {
		return
	}
	for _, result := range model.Result {
		processMessage(result)
	}
}

func processMessage(model resultModel) {
	if model.Message.Text == "/start" {
		message := fmt.Sprint("chat_id: ", model.Message.Chat.Id)
		sendMessage(model.Message.Chat.Id, message)
		return
	}

	if model.Message.Text == "/end" {
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}
	fmt.Println(model.Message.Text)
}

func sendMessage(chatID int64, text string) error {
	if chatID == 0 {
		return errors.New("invalid chat id")
	}
	postBody, err := json.Marshal(map[string]interface{}{
		"text":    text,
		"parse_mode": "markdown",
		"chat_id": chatID,
	})
	if err != nil {
		return err
	}
	responseBody := bytes.NewBuffer(postBody)
	url := telegramBaseURL() + "/sendMessage"
	_, err = http.Post(url, "application/json", responseBody)
	return err
}

func readAll() {
	nChunks := int64(0)
	r := bufio.NewReader(os.Stdin)
	buf := make([]byte, 0, 4*1024)

	for {

		n, err := r.Read(buf[:cap(buf)])
		buf = buf[:n]

		if n == 0 {
			if err == nil {
				continue
			}

			if err == io.EOF {
				break
			}

			log.Fatal(err)
		}

		nChunks++
		sendMessage(chatID, "```\n" + string(buf) + "\n```")

		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
	}
}

func main() {
	var err error
	telegramKey = os.Getenv("ttp_bot_token")
	chatID, err = strconv.ParseInt(os.Getenv("ttp_chat_id"), 10, 64)
	if err != nil {
		chatID = 0
	}
	var sendOnly = false
	var ttl = 60 //seconds
	var key string
	for _, value := range os.Args[1:] {
		if key == "" {
			key = value
			continue
		}
		switch key {
		case "--send_only":
			sendOnly, err = strconv.ParseBool(value)
			if err != nil {
				panic(err)
			}
		case "--ttl":
			ttl, err = strconv.Atoi(value)
			if err != nil {
				panic(err)
			}
		case "--chat_id":
			chatID, err = strconv.ParseInt(value, 10, 64)
			if err != nil {
				panic(err)
			}
		case "--telegramKey":
			telegramKey = value
		default:
			panic(errors.New("Invalid parameter: " + key))
		}
		key = ""
	}

	if telegramKey == "" {
		panic(errors.New("ttp_bot_token - required"))
	}


	if ttl > 0 {
		go killTheBot(ttl)
	}

	if sendOnly {
		readAll()
	} else {
		go readAll()
		runBot()
	}


}

func killTheBot(timeout int) {
	time.Sleep(time.Duration(timeout) * time.Second)
	os.Exit(0)
}

