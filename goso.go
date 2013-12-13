package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/VividCortex/robustly"
)

var sleepTime int64 = 5
var startTime int64 = time.Now().Unix() - (60 * sleepTime)

var query *string = flag.String("tags", "go", "SO query")

func main() {
	flag.Parse()

	go robustly.Run(func() { loop() })
	select {}
}

func loop() {
	for {
		qs := getLatestSOQuestions()
		for _, q := range qs {
			fmt.Println(q)
			notifyCmd := exec.Command("growlnotify", "-m", q.Title, "--url", q.Link, "--image", "so.jpeg", "SO question on "+*query)
			err := notifyCmd.Run()
			if err != nil {
				fmt.Println(err)
			}
			time.Sleep(10 * time.Second)
		}
		time.Sleep(time.Duration(sleepTime) * time.Minute)
	}
}

type SOQueryResponse struct {
	Items        []SOItem `json:"items"`
	Backoff      uint     `json:"backoff"`
	ErrorName    string   `json:"error_name"`
	ErrorMessage string   `json:"error_message"`
}

type SOItem struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

func getLatestSOQuestions() []SOItem {
	t := time.Now().Unix() - (60 * (sleepTime + 3))
	timeStr := fmt.Sprintf("%d", t)
	url := "https://api.stackexchange.com/2.1/search?fromdate=" + timeStr + "&order=asc&sort=creation&tagged=" + *query + "&site=stackoverflow"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	response := SOQueryResponse{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		fmt.Println(err)
	}
	if len(response.Items) == 0 && response.ErrorName != "" {
		fmt.Println(response.ErrorName)
		fmt.Println(response.ErrorMessage)
		if response.ErrorName == "throttle_violation" {
			split := strings.Split(response.ErrorMessage, " ")
			secs, err := strconv.Atoi(split[len(split)-2])
			if err != nil {
				fmt.Println(err)
			} else {
				if secs < 100000 {
					fmt.Println(fmt.Sprintf("throttled, sleeping for %ds", secs))
					time.Sleep(time.Duration(secs+1) * time.Second)
				}
			}
		}
	}
	if response.Backoff > 0 {
		fmt.Println(fmt.Sprintf("backoff set, sleeping for %ds", response.Backoff))
		time.Sleep(time.Duration(response.Backoff+1) * time.Second)
	}
	return response.Items
}
