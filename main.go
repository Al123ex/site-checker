package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const TOKEN string = ""
const CHANNEL string = ""
const DELAY time.Duration = 5
const FILENAME string = "sites.txt"

type site struct {
	Url    string
	Status bool
}

func main() {
	bytes, _ := ioutil.ReadFile(FILENAME)
	links := strings.Split(string(bytes), "\n")

	c := make(chan site)

	for _, link := range links {
		site := site{link, true}
		go checkLink(site, c)
	}

	for s := range c {
		go func(s site) {
			time.Sleep(DELAY * time.Second)
			checkLink(s, c)
		}(s)
	}
}

func checkLink(s site, c chan site) {
	_, err := http.Get(s.Url)
	if err != nil {
		if s.Status {
			message := "<b>Отчет за </b>" + time.Now().Format("2.1 15:04")
			message += "%0A"
			message += "Сайт <a href=\"" + s.Url + "\">" + s.Url + "</a> не работает."
			send(message)
		}
		s.Status = false
	} else {
		s.Status = true
	}

	c <- s
}

func send(message string) bool {
	link := "https://api.telegram.org/bot" + TOKEN + "/sendMessage?chat_id=" + CHANNEL
	link += "&text=" + message + "&parse_mode=HTML"

	_, err := http.Get(link)

	if err != nil {
		fmt.Println(err.Error())
	}

	return err != nil
}
