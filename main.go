package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"sort"
	"strings"
)

var client http.Client

type Leaderboard struct {
	OwnerId int               `json:"owner_id"`
	Event   string            `json:"event"`
	Members map[string]Member `json:"members"`
}

type Member struct {
	LocalScore int    `json:"local_score"`
	Stars      int    `json:"stars"`
	Name       string `json:"name"`
}

func init() {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Error while creating cookie jar %s", err)
	}
	client = http.Client{
		Jar: jar,
	}
}

type Members []Member

func (s Members) Len() int {
	return len(s)
}

func (s Members) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Members) Less(i, j int) bool {
	return s[i].LocalScore > s[j].LocalScore
}

func main() {
	cookieFile := flag.String("f", ".cookie", "a path to the cookie")
	leaderBoardJSON := flag.String("l", "", "the URL for the leaderboard json")
	flag.Parse()
	if *leaderBoardJSON == "" {
		fmt.Println("Please provide a leaderboard link!")
		os.Exit(1)
	}
	cookieVal, err := os.ReadFile(*cookieFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	cookieClean := strings.TrimSuffix(string(cookieVal), "\n")
	cookie := &http.Cookie{
		Name:  "session",
		Value: cookieClean,
	}
	// https://adventofcode.com/2022/leaderboard/private/view/1699975.json
	req, err := http.NewRequest("GET", *leaderBoardJSON, nil)
	if err != nil {
		log.Fatal("Error with request")
	}
	req.AddCookie(cookie)

	resp, _ := client.Do(req)
	defer resp.Body.Close()
	resBody, _ := io.ReadAll(resp.Body)
	leaderboard := Leaderboard{}
	json.Unmarshal(resBody, &leaderboard)
	members := leaderboard.Members
	var people []Member
	for _, val := range members {
		people = append(people, val)
	}
	sort.Sort(Members(people))
	for _, person := range people {
		stars := strings.Repeat("⭐", person.Stars)
		fmt.Printf("%s %d %s\n", person.Name, person.LocalScore, stars)
	}
}
