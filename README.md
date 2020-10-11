# Systems Assignment

## What is it?

This exercise is a follow-on to the [General Assignment](https://github.com/cloudflare-hiring/cloudflare-2020-general-engineering-assignment), you'll need to complete that first.  In this assignment you'll write a program that makes a request to the endpoints you created in the General Assignment.  This is a systems assignment so we want to see that you're able to use sockets directly rather than using a library that handles the HTTP request.

## Implementation Details

- Written in Go
- Implemented via sockets to perform HTTP GET request
- Keep-alive is not used. Connection was hanging for some reason when after receiving the header and waiting for the body

## To build
- Run `make`

## To run
- Run `make run` to make a request to `http://cloudflare-workers.akhtariev.workers.dev:80/links`
- Run `systems-assignment` executable with the following flags:

`--url`       String  full URL. Must include port. e.g. https://example.org:80/path
`--profile`   Int     positive integer for the number of requests to profile

Possible flag combinations:
`--url`                   performs HTTP GET request.")
`--url` and `--profile`   performs HTTP GET request and profiles the page with number of requests equal to profile.
`--help`                  to print the above information about the flags

### Note on `--profile`

- Measures the time elapsed from before the beginning of the write to socket and up until the receival of the first token (as per Scanner.Scan)

## Screenshots

##### Cloudflare
- `http://cloudflare-workers.akhtariev.workers.dev:80/links`

  ![](/screenshots/cloudflare.PNG)

  

##### Other
- `http://jsonplaceholder.typicode:80.com/todos/1`

![](/screenshots/other.PNG)



## Observations ##

Cloudflare's response is 6 times bigger but it is approximately as performant as the other server. Moreover, it is faster in the best case