package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type count struct {
	right int
	total int
}

// Read unser from console and return it
func consoleReader() string {
	r := bufio.NewReader(os.Stdin)
	un, err := r.ReadString('\n')

	if err != nil {
		log.Fatalln("Input error", err)
	}

	return un
}

func main() {
	ask := true

	// Parse flags
	var path = flag.String("file", "problems.csv", "")
	var timeLimit = flag.String("time", "10", "")
	flag.Parse()

	// Prepare result
	c := count{right: 0, total: 0}

	// Read csv file
	csvFile, err := os.Open(*path)

	if err != nil {
		log.Fatalln("File error", err)
	}

	r := csv.NewReader(csvFile)

	// Start timer
	timeLimitInt, err := strconv.Atoi(*timeLimit)
	if err != nil {
		log.Fatalln("Must be a number in seconds", err)
	}

	timeout := time.After(time.Duration(timeLimitInt) * time.Second)

	// Base
	for {
		record, err := r.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalln("Can't read file", err)
		}

		select {

		// Wait for timeout
		case <-timeout:
			c.total++
			ask = false
			break

		// Ask question and check answer. Last answer will be checked, even if timout comes
		default:
			if ask {
				fmt.Print(record[0], " =", " ")
				correctAnswer := record[1]
				answer := strings.TrimRight(consoleReader(), "\n")
				fmt.Println("here", answer)
				if answer == correctAnswer {
					c.right++
				}
			}

			c.total++
		}
	}
	fmt.Println("Correct unswers: ", c.right, "Total unswers: ", c.total)
}
