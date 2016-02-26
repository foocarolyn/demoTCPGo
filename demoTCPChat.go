package main

import ("log"; "net"; "os"; "fmt";)

type Client struct {
	conn net.Conn
	nickname string
	msgChannel chan string
}

func waitForNickname(conn net.Conn) string {
	buf := make([]byte, 4)

	// Read nick name size
	n, err := conn.Read(buf)
	if err != nil || n != 4 {
		log.Printf("err reading nickname size: %s, n: %d", err, n)
	}

	nicknameSize := Byte32bArrToInt(buf)
	nicknameBuff := make([]byte, nicknameSize)
	n, err = conn.Read(nicknameBuff)
	if err != nil || n != int(nicknameSize) {
		log.Printf("err reading nickname: %s, readSize: %d intendedSize:%d", err, n, nicknameSize)
		return ""
	}
	return string(nicknameBuff[:n])
}

func handleConnection(conn net.Conn, centralMsgChannel chan <- string, addClientChannel chan <- Client, rmClientChannel chan <- Client) {
	defer conn.Close()

	nickname := waitForNickname(conn)
	client := Client {
		conn: conn,
		nickname: nickname,
		msgChannel: make(chan string),
	}
	log.Printf("Client: %s joined", nickname)
	// addClientChannel <- client

	defer func () {
		// rmClientChannel <- client
		// centralMsgChannel <- fmt.Sprintf("User %s has left the chat room \n", client.nickname)
		// log.Printf("Connection from %v closed. \n", conn.RemoteAddr())
	}()

	welcomeString := fmt.Sprintf("Welcome, %s", client.nickname)

	headerLength := 1 /*cmd*/
	totalLength := headerLength + len(welcomeString)
	headerBuf := IntTo32bByteArr(totalLength)

	n, error := conn.Write(headerBuf)
	n, error = conn.Write(append([]byte{0x1}, []byte(welcomeString)...))

	if error != nil {
		log.Printf("%d, %s", n, error.Error())
	}

	// centralMsgChannel <- fmt.Sprintf("New user %s has joined the chat room \n", client.nickname)

	// go client.ReadLinesInto(centralMsgChannel)
	// client.WriteLines()
}

func main() {
	ln, err := net.Listen("tcp", ":6000")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	log.Printf("Listening to tcp 6000")

	// Make the channels
	centralMsgChannel := make(chan string)
	addClientChannel := make(chan Client)
	rmClientChannel := make(chan Client)

	// go handleMessages(centralMsgChannel, addClientChannel, rmClientChannel)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}

		go handleConnection(conn, centralMsgChannel, addClientChannel, rmClientChannel)
	}
}