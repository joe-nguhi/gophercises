package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

type Problem struct {
	Question string
	Answer   string
}

var count int

func main() {

	problems := getProblems("problems.csv")
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Starting Quiz....")
	fmt.Println("Press Enter key to move to next Quiz Question. Good Luck!")

	for i, problem := range problems {
		fmt.Printf("%d. %s = ", i+1, problem.Question)

		scanner.Scan()
		ans := scanner.Text()

		if ans == problem.Answer {
			count++
		}
	}

	fmt.Printf("You scored %d out of %d\n", count, len(problems))
}

func getProblems(filepath string) []*Problem {
	file, err := os.Open(filepath)

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
