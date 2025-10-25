package main

import (
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

func main() {

	problems := getProblems("problems.csv")

	for i, problem := range problems {
		fmt.Printf("%d. %s = %s\n", i+1, problem.Question, problem.Answer)
	}
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
