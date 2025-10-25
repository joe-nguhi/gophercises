package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type Problem struct {
	Question string
	Answer   string
}

const defaultLimit = 15
const defaultCSV = "problems.csv"

var count int

var flagLimit int
var flagCSV string

func init() {
	flag.StringVar(&flagCSV, "csv", defaultCSV, "a csv file in the format 'question,answer'")
	flag.IntVar(&flagLimit, "limit", defaultLimit, "the time limit for the quiz in seconds")
	flag.Parse()
}

func main() {
	problems := getProblems(flagCSV)
	timer := time.After(time.Duration(flagLimit) * time.Second)

	completed := make(chan bool)

	go runQuiz(problems, completed)

	for {
		select {
		case <-timer:
			fmt.Println("\nTime's up!")
			fmt.Printf("You scored %d out of %d\n", count, len(problems))
			os.Exit(0)
		case <-completed:
			fmt.Printf("You scored %d out of %d\n", count, len(problems))
			os.Exit(0)
		default:
		}
	}
}

func getProblems(filepath string) []*Problem {
	file, err := os.Open(filepath)
	defer file.Close()

	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(file)

	var problems []*Problem

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		problems = append(problems, &Problem{
			record[0], record[1],
		})
	}

	return problems
}

func runQuiz(problems []*Problem, done chan bool) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Starting Quiz....")
	fmt.Println("Press Enter key to move to next Quiz Question. Good Luck!")

	for i, problem := range problems {
		fmt.Printf("%d. %s = ", i+1, problem.Question)

		scanner.Scan()
		ans := scanner.Text()

		parsed := strings.ToLower(strings.TrimSpace(ans))

		if parsed == strings.ToLower(problem.Answer) {
			count++
		}
	}

	done <- true
}
