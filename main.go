package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/c-bata/go-prompt"
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
		p := prompt.New(func(input string) {
			if conn == nil {
				println("No Client")
				return
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			if err := conn.Write(ctx, websocket.MessageText, []byte(input)); err != nil {
				log.Println(err)
			}

			cancel()
		}, func(input prompt.Document) []prompt.Suggest {
			return []prompt.Suggest{}
		}, prompt.OptionPrefix(">>> "), prompt.OptionTitle("DebugWS"))

		defer os.Exit(0)
		p.Run()
	}()

	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Fatal(err)
	}
}
