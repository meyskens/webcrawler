A simple (~~not~~ now concurrent) web crawler
=====================================

*This code is made for my coding challenge when applying for an internship at GoCardless, later made concurrent*

## Dependencies
This app is written in Go and isn't compiles so you need to have Go installed. You can get it at https://golang.org/dl/

This app doesn't use any dependencies apart from the standard library. So `go get` isn't needed , in case of any build error it is suggested to run it.

## Building
To build the binary you can use the command `go build ./` which will make a binary with the name of the directory the file is in. 

## Run
### without build
You can also run the code by using `go run main.go "https://gocardless.com"`, the first argument after the file name is the URL to crawl
### with build
Just run the program in your favorite shell, note that the filename depends on the directory you compiled it in. `./crawler "https://gocardless.com"` , again the second argument is the URL