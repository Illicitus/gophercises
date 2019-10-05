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
)

type count struct {
	right int
	total int
}

// Read anser from console and return it
func consoleReader() string {
	r := bufio.NewReader(os.Stdin)
	an, err := r.ReadString('\n')

	if err != nil {
		log.Fatalln("Input error", err)
	}

	return an
}

func main() {

	// Parse flag
	var path = flag.String("file", "problems.csv", "")
	flag.Parse()

	// Prepare result
	c := count{right: 0, total: 0}

	// Read csv file
	csvFile, err := os.Open(*path)

	if err != nil {
		log.Fatalln("File error", err)
	}

	r := csv.NewReader(csvFile)

	for {
		record, err := r.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalln("Can't read file", err)
		}

		// Aks question and check answer
		fmt.Print(record[0], " =", " ")
		cr := record[1]
		an := strings.TrimRight(consoleReader(), "\n")

		if an == cr {
			c.right++
		}

		c.total++

	}
	fmt.Println("Correct answers: ", c.right, "Total answers: ", c.total)
}
