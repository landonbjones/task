package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strconv"
)

const nameEndpoint = "https://names.mcquay.me/api/v0"
const jokeEndpoint = "http://api.icndb.com/jokes/random"

type Name struct {
	First string `json:"first_name"`
	Last  string `json:"last_name"`
}

var port = flag.Int("port", 5000, "the port number for the web server")

func main() {
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		name, err := randomName()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		joke, err := randomJoke(name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, joke)
	})

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}

func randomName() (Name, error) {
	var name Name
	resp, err := http.Get(nameEndpoint)
	if err != nil {
		return name, logError(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return name, logError(err)
	}

	err = json.Unmarshal(body, &name)
	return name, logError(err)
}

func randomJoke(name Name) (string, error) {
	resp, err := http.Get(jokeEndpoint + "?limitTo=nerdy&firstName=" + name.First + "&lastName=" + name.Last)
	if err != nil {
		return "", logError(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", logError(err)
	}

	type Value struct {
		Id   int    `json:"id"`
		Joke string `json:"joke"`
	}

	type JokeResponse struct {
		Type  string `json:"type"`
		Value Value  `json:"value"`
	}

	var jokeResponse JokeResponse
	err = json.Unmarshal(body, &jokeResponse)
	return jokeResponse.Value.Joke, logError(err)
}

func logError(err error) error {
	if err == nil {
		return nil
	}
	_, filename, line, _ := runtime.Caller(1)
	log.Printf("[error] %s:%d %v", filename, line, err)
	return err
}
