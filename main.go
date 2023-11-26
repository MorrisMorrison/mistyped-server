package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

var keyMap = map[string]string{
	"q": "aw", "w": "se", "e": "dr", "r": "ft", "t": "gz", "z": "hu", "u": "ji", "i": "ko", "o": "lp", "p": "öü", "ü": "ä",
	"a": "ys", "s": "xd", "d": "cf", "f": "vg", "g": "bh", "h": "nj", "j": "mk", "k": "l", "l": "ö", "ö": "ä", "ä": "",
	"y": "sx", "x": "dc", "c": "fv", "v": "gb", "b": "hn", "n": "jm", "m": "k",
}

var port = 8080

type CheckUrlResponse struct {
	Available   []string `json:"Available"`
	Unavailable []string `json:"Unavailable"`
}

func main() {
	handleRequests()
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/checkurl/", handleCheckUrl)
	portString := ":" + strconv.Itoa(port)
	fmt.Println("mistyped-server up and running on port", portString)
	log.Fatal(http.ListenAndServe(portString, router))
}

func handleCheckUrl(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ENDPOINT CHECKURL")
	fmt.Println("Request:", r)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST, PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	inputUrl := r.Form.Get("url")
	if !isUrlValid(inputUrl) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cleanInputUrl := getCleanUrl(inputUrl)
	candidates := getCandidates(cleanInputUrl)

	c1, c2 := checkUrlAvailability(candidates)
	availableCandidates := <-c1
	unavailableCandidates := <-c2

	fmt.Println(availableCandidates)
	fmt.Println(unavailableCandidates)

	response := CheckUrlResponse{
		Available:   availableCandidates,
		Unavailable: unavailableCandidates,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func checkUrlAvailability(candidates []string) (chan []string, chan []string) {
	c1 := make(chan []string)
	c2 := make(chan []string)

	go func() {
		fmt.Println("START GOROUTINE CheckUrlAvailability")
		availableCandidates := make([]string, 0)
		unavailableCandidates := make([]string, 0)

		for _, candidate := range candidates {
			client := http.Client{Timeout: 1 * time.Second}
			res, err := client.Get(candidate)
			if err == nil && res != nil && (res.StatusCode == 200 || res.StatusCode == 204 || res.StatusCode == 403) {
				availableCandidates = append(availableCandidates, candidate)
			} else {
				unavailableCandidates = append(unavailableCandidates, candidate)
			}
		}

		c1 <- availableCandidates
		c2 <- unavailableCandidates

		close(c1)
		close(c2)
	}()

	return c1, c2
}

func getCandidates(inputUrl string) []string {
	candidates := make([]string, 0)
	urlSplit := strings.Split(inputUrl, ".")
	var hostPosition int
	if len(urlSplit) > 2 {
		hostPosition = 1
	} else {
		hostPosition = 0
	}
	hostname := urlSplit[hostPosition]

	for _, urlCharacter := range urlSplit[hostPosition] {
		for _, possibleCharacters := range keyMap[string(urlCharacter)] {
			for _, possibleCharacter := range string(possibleCharacters) {
				possibleHostName := replace(hostname, string(urlCharacter), string(possibleCharacter))
				urlSplit[hostPosition] = possibleHostName
				possibleUrl := getString(urlSplit)
				candidates = append(candidates, "http://"+possibleUrl)
			}
		}
	}

	return candidates
}

func replace(replaceIn string, replaceThis string, replaceWith string) string {
	result := ""
	for _, c := range replaceIn {
		if string(c) == replaceThis {
			result += replaceWith
		} else {
			result += string(c)
		}
	}
	return result
}

func getString(stringArray []string) string {
	result := ""
	for i, s := range stringArray {
		result += s

		if i != len(stringArray)-1 {
			result += "."
		}
	}

	return result
}

func getCleanUrl(inputUrl string) string {
	url := strings.TrimPrefix(inputUrl, "https://")
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "www.")
	url = strings.Split(url, "/")[0]
	url = strings.TrimSuffix(url, "/")
	return url
}

func isUrlValid(u string) bool {
	var err error

	if len(strings.Split(u, ".")) <= 2 {
		_, err = url.Parse("http://www." + u)
	} else {
		_, err = url.Parse(u)
	}

	fmt.Println(err)
	return err == nil
}
