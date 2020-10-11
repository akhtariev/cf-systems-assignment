package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"math"
	"net"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

type durations []int64

func (durations durations) Len() int           { return len(durations) }
func (durations durations) Less(i, j int) bool { return durations[i] < durations[j] }
func (durations durations) Swap(i, j int)      { durations[i], durations[j] = durations[j], durations[i] }

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func validateFlags(help bool, url string, profile int) (bool, error) {
	if flag.NFlag() > 2 || flag.NFlag() == 0 {
		return true, errors.New("Invalid number of flags. Add --help to see available flags")
	} else if flag.NFlag() == 2 {
		if url == "" || profile < 0 {
			return true, errors.New("Invalid flags supplied. Add --help to see available flags")
		}
	} else if help == true {
		fmt.Println("Available flags:")
		fmt.Println("--url   String  Full URL. e.g. https://example.org:8000/path")
		fmt.Print("--profile   Int     positive integer for the number of requests to profile\n\n")
		fmt.Println("Possible flag combinations:")
		fmt.Println("--url             performs HTTP GET request.")
		fmt.Println("--url --profile   performs HTTP GET request and profiles the page with number of requests equal to profile. ")
	} else if url == "" {
		return true, errors.New("Invalid flag supplied. Add --help to see available flags")
	}
	return false, nil
}

func appendSlice(slice []string, data ...string) []string {
	m := len(slice)
	n := m + len(data)
	if n > cap(slice) { // if necessary, reallocate
		// allocate double what's needed, for future growth.
		newSlice := make([]string, (n+1)*2)
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[0:n]
	copy(slice[m:n], data)
	return slice
}

func main() {
	helpPtr := flag.Bool("help", false, "Help command")
	urlPtr := flag.String("url", "", "URL")
	profilePtr := flag.Int("profile", -1, "URL")
	flag.Parse()

	shouldExit, err := validateFlags(*helpPtr, *urlPtr, *profilePtr)
	checkError(err)
	if shouldExit {
		os.Exit(0)
	}

	u, err := url.Parse(*urlPtr)
	checkError(err)

	var path string
	if u.EscapedPath() == "" {
		path = "/"
	} else {
		path = u.EscapedPath()
	}

	var profileCount int
	if *profilePtr > 0 {
		profileCount = *profilePtr
	} else {
		profileCount = 1
	}

	var start time.Time
	var end time.Time
	var currentDurationMs int64
	var sumDurationMs int64 = 0
	var minDuration int64 = math.MaxInt64
	var maxDuration int64 = 0
	var durations durations = make(durations, profileCount)
	var errorCodes []string = []string{}
	var currentErrorCode string
	reErrorCode := regexp.MustCompile(`\d\d\d`)
	reContentLength := regexp.MustCompile(`\d+`)
	var minLength string
	var maxLength string

	for i := 0; i < profileCount; i++ {
		conn, err := net.Dial("tcp", u.Hostname()+":"+u.Port())
		checkError(err)
		start = time.Now()
		_, err = conn.Write([]byte("GET " + path + " HTTP/1.1\r\nHost: " + u.Hostname() + "\r\nAccept: application/json\r\nConnection: close\r\n\r\n"))
		checkError(err)

		scanner := bufio.NewScanner(conn)
		if scanner == nil || scanner.Scan() == false {
			fmt.Println("Problems receiving response from the server.")
			os.Exit(1)
		}
		end = time.Now()

		statusLine := scanner.Text()
		if !strings.HasSuffix(statusLine, "200 OK") {
			currentErrorCode = string(reErrorCode.Find([]byte(statusLine)))
			fmt.Println("Server failure. Abort with code: " + currentErrorCode)
			errorCodes = appendSlice(errorCodes, currentErrorCode)
			continue
		}

		currentDurationMs = end.Sub(start).Milliseconds()
		sumDurationMs += currentDurationMs
		if currentDurationMs > maxDuration {
			maxDuration = currentDurationMs
		}

		if currentDurationMs < minDuration {
			minDuration = currentDurationMs
		}
		durations[i] = currentDurationMs

		// Scan over the header until CRLF (according to RFC)
		for scanner.Scan() {
			if scanner.Text() == "" {
				break
			} else if strings.HasPrefix(scanner.Text(), "Content-Length:") {
				if currentDurationMs == maxDuration {
					maxLength = string(reContentLength.Find([]byte(scanner.Text())))
				}

				if currentDurationMs == minDuration {
					minLength = string(reContentLength.Find([]byte(scanner.Text())))
				}
			}
		}

		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		conn.Close()
		fmt.Println("")
	}

	if *profilePtr > 0 {
		sort.Sort(durations)
		var median int64
		if len(durations)%2 == 0 {
			median = (durations[profileCount/2] - durations[profileCount/2-1]) / 2
		} else {
			median = durations[profileCount/2]
		}

		fmt.Println("\nProfile Information:")
		fmt.Printf("Number of requests: %d\n", profileCount)
		fmt.Printf("Fastest time: %d ms\n", minDuration)
		fmt.Printf("Slowest time: %d ms\n", maxDuration)
		fmt.Printf("Mean time: %d ms\n", sumDurationMs/int64(profileCount))
		fmt.Printf("Median time: %d ms\n", median)
		fmt.Print("Error Codes: ")
		if len(errorCodes) > 0 {
			fmt.Print("\n")
			for i := 0; i < len(errorCodes); i++ {
				fmt.Println("- " + errorCodes[i])
			}
		} else {
			fmt.Println("None")
		}
		fmt.Printf("Smallest response size: %s bytes\n", minLength)
		fmt.Printf("Largest response size: %s bytes\n", maxLength)
	}

	os.Exit(0)
}
