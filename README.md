# **Talker Chat**

Talker Chat is an applications that allow multiple clients talk through a chat server. It's implemented using the TCP protocol and also make use of some important pattern in Golang.

## **Prerequisites**

This application was developed using the version 1.9 of Go, so be aware.

## **Dependencies**

There aren't any dependencies. This application make use only of the standard libs of Go.

## **Running**

* First, clone this repo inside $GOPATH/src/github.com
`$GOPATH` is the environment variable that points to your Go workspace.

* Access the *wand-go-talker-chat/execution* folder through terminal and execute `go run chatServer.go`.

* Now, open at least two others terminal sessions and *in both of them* access the *wand-go-talker-chat/chatclient* folder and execute `go run chat client.go`.

* Follow what is asked and see the magic happening.

## License

Copyright © 2017 Wanderson Silva