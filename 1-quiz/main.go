package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"
)

var problemFile string
var timeLimit int

func init() {
	// Parse command line flags, use default if not passed
	flag.StringVar(&problemFile, "csv", "problems.csv", "a csv file with questions and answers")
	flag.IntVar(&timeLimit, "limit", 30, "the time limit for the quiz in seconds")

	flag.Parse()
}

func main() {
	// Open the csv file for reading
	file, err := os.Open(problemFile)
	if err != nil {
		log.Fatal(err)
	}
	// Defer closing the file after it is not used anymore
	defer file.Close()

	// Create a reader for parsing the csv file
	csv := csv.NewReader(file)

	// Parse the csv file into a mapped string
	problems, err := csv.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	total, correct, wrong := len(problems),0,0

	// Create a timer with the defined limit
	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)

	// Wait in a new thread until the timer is finished
	go func() {
		<-timer.C

		// Print out the total score after finishing the test
		fmt.Printf("\nYou scored %d out of %d (%d%%)\n", correct, total, 100.0*correct/total)

		// Exit the application
		os.Exit(0)
	}()

	// Loop through the csv file, reading each line
	for i, problem := range problems {
		// After the last line was read, exit the loop
		if err == io.EOF {
			break
		}

		fmt.Printf("Problem #%d: %s = ", i+1, problem[0])

		// Get the user input
		var input string
		fmt.Scanln(&input)

		// Check the input against the answer
		if input != problem[1] {
			wrong++
		} else {
			correct++
		}

		i++
	}

	// Print out the total score after finishing the test
	fmt.Printf("You scored %d out of %d (%d%%)\n", correct, total, 100.0*correct/total)
}
