package main

import (
    "net"
    "fmt"
    "os"
    "flag"
    "bufio"
    "strings"
    "regexp"
    "net/url"
    "errors"
)

func checkError(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
        os.Exit(1)
    }
}

func validateFlags(help bool, url string, profile int) (bool, error) {
    if (flag.NFlag() > 2 || flag.NFlag() == 0) {
        return true, errors.New("Invalid number of flags. Add --help to see available flags");
    } else if (flag.NFlag() == 2) {
        if (url == "" || profile < 0) {
            return true, errors.New("Invalid flags supplied. Add --help to see available flags")
        }
    } else if (help == true) {
        fmt.Println("Available flags:")
        fmt.Println("--url   String  Full URL. e.g. https://example.org:8000/path")
        fmt.Print("--profile   Int     positive integer for the number of requests to profile\n\n")
        fmt.Println("Possible flag combinations:")
        fmt.Println("--url             performs HTTP GET request.")
        fmt.Println("--url --profile   performs HTTP GET request and profiles the page with number of requests equal to profile. ")
    }  else if (url == "") {
        return true, errors.New("Invalid flag supplied. Add --help to see available flags")
    }
    return false, nil
}

func main() {
    helpPtr := flag.Bool("help", false, "Help command")
    urlPtr := flag.String("url", "", "URL")
    profilePtr := flag.Int("", -1, "URL")
    flag.Parse()

    shouldExit, err := validateFlags(*helpPtr, *urlPtr, *profilePtr)
    checkError(err)
    if (shouldExit) {
        os.Exit(0)
    }

    u, err := url.Parse(*urlPtr)
    checkError(err)

    conn, err := net.Dial("tcp", u.Hostname() + ":" + u.Port())
    checkError(err)
  
    var path string;
    if (u.EscapedPath() == "") {
        path = "/"
    } else {
        path = u.EscapedPath()
    }
    _, err = conn.Write([]byte("GET " + path + " HTTP/1.0\r\nHost: " + u.Hostname() + "\r\n\r\n"))
    checkError(err)
  
    scanner := bufio.NewScanner(conn)
    if (scanner == nil || scanner.Scan() == false) {
        fmt.Println("Problems receiving response from the server.")
        os.Exit(1)
    }

    statusLine := scanner.Text()
    if(!strings.HasSuffix(statusLine, "200 OK")) {
        re := regexp.MustCompile(`\d\d\d`)
        fmt.Println("Server failure. Abort with code: " + string(re.Find([] byte(statusLine))))
        os.Exit(1)
    }

    // Scan over the header until CRLF (according to RFC)
    for (scanner.Scan()) {
        if (scanner.Text() == "") {
            break
        }
    }

    for (scanner.Scan()) {
        fmt.Println(scanner.Text())
    }
    os.Exit(0)
}
