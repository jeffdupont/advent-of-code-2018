package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	re := regexp.MustCompile("[0-9]+")

	var closetLog []entry

	s := bufio.NewScanner(f)
	for s.Scan() {
		var timestamp, action string

		line := s.Text()
		segments := strings.Split(line, " ")

		timestamp = strings.Trim(segments[0]+" "+segments[1], "[]")

		t, err := time.Parse("2006-01-02 15:04", timestamp)
		if err != nil {
			log.Fatal(err)
		}

		action = strings.Join(segments[2:], " ")
		id, _ := strconv.Atoi(re.FindString(action))
		closetLog = append(closetLog, entry{timestamp: t, action: action, id: id})
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	sort.Slice(closetLog, func(i, j int) bool {
		return closetLog[i].timestamp.Before(closetLog[j].timestamp)
	})

	sleeper := map[int]int{}
	guard := 0
	for i := 0; i < len(closetLog); i++ {
		if guard != closetLog[i].id && closetLog[i].id != 0 {
			guard = closetLog[i].id
		}
		fmt.Println(i, closetLog[i].timestamp.Format("2006-01-02 15:04"), closetLog[i].action)
		if closetLog[i].action == "falls asleep" {
			mins := minutesAsleep(closetLog[i].timestamp, closetLog[i+1].timestamp)
			fmt.Println(guard, mins)
			sleeper[guard] += mins
			i++
		}
	}

	fmt.Println(sleeper)
}

func minutesAsleep(a, b time.Time) int {
	return int(b.Sub(a).Minutes())
}

type entry struct {
	timestamp time.Time
	action    string
	id        int
}
