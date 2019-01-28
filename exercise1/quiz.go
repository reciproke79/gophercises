package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type problem struct {
	q, a string
}

func main() {
	fileFlag := flag.String("csv", "problems.csv", "a csv file in the format 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit of the quiz in seconds")
	flag.Parse()

	f, err := os.Open(*fileFlag)
	if err != nil {
		exit(fmt.Sprintf("Can't open file: %s", *fileFlag))
	}

	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		exit(fmt.Sprintf("Error reading CSV file, please check format."))
	}

	problems := parseLines(records)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.q)

		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("\nYou answered %d out of %d questions correctly\n", correct, len(problems))
			return
		case answer := <-answerCh:
			if answer == p.a {
				correct++
			}
		}
	}
	fmt.Printf("\nYou answered %d out of %d questions correctly\n", correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	var problems []problem
	for _, line := range lines {
		problems = append(problems, problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		})
	}
	return problems
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
