package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type data struct {
	Level   int
	Message string
}

var conn *websocket.Conn

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		c, err := websocket.Accept(w, req, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer c.Close(websocket.StatusInternalError, "the sky is falling")
		c.SetReadLimit(8192000)
		conn = c

		log.Println("Debug ws connected " + req.RemoteAddr)

		for {
			var v data
			if err = wsjson.Read(req.Context(), c, &v); err != nil {
				log.Println(err)
				conn = nil
				return
			}

			switch v.Level {
			case 0:
				log.Print("T: ")
			case 1:
				log.Print("I: ")
			case 2:
				log.Print("W: ")
			case 3:
				log.Print("E: ")
			}
			fmt.Println(v.Message)
		}
	})

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for {
			if conn != nil {
				scanner.Scan()
				code := scanner.Text()

				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

				if err := conn.Write(ctx, websocket.MessageText, []byte(code)); err != nil {
					log.Println(err)
				}
				cancel()
			}
		}
	}()

	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Fatal(err)
	}
}
