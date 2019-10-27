package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

// Complete the camelcase function below.
func camelcase(s string) int32 {
	count := 0

	if len(s) > 0 {
		count++
	}

	for _, i := range s {
		if unicode.IsUpper(i) {
			count++
		}
	}

	return int32(count)
}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1024*1024)

	s := readLine(reader)

	result := camelcase(s)

	fmt.Println(result)
}

func readLine(reader *bufio.Reader) string {
	str, _, err := reader.ReadLine()
	if err == io.EOF {
		return ""
	}

	return strings.TrimRight(string(str), "\r\n")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
