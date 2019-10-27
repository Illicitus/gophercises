package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Complete the caesarCipher function below.
func caesarCipher(s string, k int32) string {
	sm := "abcdefghijklmnopqrstuvwxyz"
	lg := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	k = indexNormilize(k)
	fmt.Println(k)

	nsm := sm[k:] + sm[:k]
	nlg := lg[k:] + lg[:k]

	var result string
	for _, v := range s {
		if i := strings.Index(sm, string(v)); i != -1 {
			result += string(nsm[i])
		} else if i := strings.Index(lg, string(v)); i != -1 {
			result += string(nlg[i])
		} else {
			result += string(v)
		}
	}

	return strings.TrimRight(result, "\r\n")
}

func indexNormilize(k int32) int32 {
	if k >= 26 {
		result := k - 26
		if result >= 26 {
			return indexNormilize(result)
		} else {
			return result
		}
	}
	return k
}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1024*1024)

	_, err := strconv.ParseInt(readLine(reader), 10, 64)
	checkError(err)

	s := readLine(reader)

	kTemp, err := strconv.ParseInt(readLine(reader), 10, 64)
	checkError(err)
	k := int32(kTemp)

	result := caesarCipher(s, k)

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
