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

	closet := closet{}

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
		closet.addEntry(entry{timestamp: t, action: action, id: id})
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	closet.sort()
	closet.processSleep()

	fmt.Println(closet.frequent())
}

type closet struct {
	log      []entry
	sleeper  map[int]int
	schedule map[int]map[int]int
}

func (c closet) frequent() int {
	guard := 0
	idealMin := 0
	min := 0
	for g := range c.sleeper {
		for m, count := range c.schedule[g] {
			if count > idealMin {
				min = m
				idealMin = count
				guard = g
			}
		}
	}
	return guard * min
}

func (c closet) okGo() int {
	maxSleep := 0
	guard := 0
	for g, sleep := range c.sleeper {
		if sleep > maxSleep {
			guard = g
			maxSleep = sleep
		}
	}
	idealMin := 0
	min := 0
	for m, count := range c.schedule[guard] {
		if count > idealMin {
			min = m
			idealMin = count
		}
	}
	return guard * min
}

func (c *closet) processSleep() {
	c.sleeper = map[int]int{}
	guard := 0
	for i := 0; i < len(c.log); i++ {
		if guard != c.log[i].id && c.log[i].id != 0 {
			guard = c.log[i].id
		}
		if c.log[i].action == "falls asleep" {
			mins := c.minutesAsleep(guard, c.log[i].timestamp, c.log[i+1].timestamp)
			c.sleeper[guard] += mins
			i++
		}
	}
}

func (c *closet) sort() {
	sort.Slice(c.log, func(i, j int) bool {
		return c.log[i].timestamp.Before(c.log[j].timestamp)
	})
}

func (c *closet) addEntry(e entry) {
	c.log = append(c.log, e)
}

func (c *closet) minutesAsleep(guard int, a, b time.Time) int {
	count := 0
	if c.schedule == nil {
		c.schedule = make(map[int]map[int]int)
	}
	for a.Before(b) {
		count++
		if _, ok := c.schedule[guard]; !ok {
			c.schedule[guard] = make(map[int]int)
		}
		c.schedule[guard][a.Minute()]++
		a = a.Add(1 * time.Minute)
	}
	return count
}

type entry struct {
	timestamp time.Time
	action    string
	id        int
}
