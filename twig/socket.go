package twig

import (
	"encoding/json"
	"leaf/utils"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

type Twig struct {
	WS    *websocket.Conn
	Token string
}

func Connect(token string, listeningTo []string) *Twig {
	client := websocket.Dialer{}

	conn, _, err := client.Dial("ws://100.105.87.12:4000/ws", nil)
	if err != nil {
		panic(utils.ParsedError(err, "Could not connect to twig", true).Error())
	}

	println("Connected to twig")

	_, message, err := conn.ReadMessage()
	if err != nil {
		utils.ParsedError(err, "Could not read message from twig", true)
		return nil
	}

	// first message should always be the hola message.
	var holaMessage HolaMessage
	if err := json.Unmarshal(message, &holaMessage); err != nil {
		utils.ParsedError(err, "Could not parse hola response", true)
		return nil
	}

	if holaMessage.Op != 1 {
		utils.ParsedError(nil, "Expected hola response, got something else", true)
		return nil
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	// lets auth ourselves now
	if err := conn.WriteJSON(AuthMessage{
		Op: 2,
		Data: struct {
			Token       string   "json:\"token\""
			Hostname    string   "json:\"hostname\""
			ListeningTo []string "json:\"listening_to\""
		}{
			Token:       token,
			Hostname:    hostname,
			ListeningTo: listeningTo,
		},
	}); err != nil {
		utils.ParsedError(err, "Could not auth to twig", true)
	}

	println("Authenticated to twig")

	go readMessage(*conn)
	go heartbeat(*conn, time.Duration(holaMessage.Data.HeartbeatInterval)*time.Millisecond)

	return &Twig{
		Token: token,
		WS:    conn,
	}
}

func readMessage(conn websocket.Conn) {
	for {
		_, _, err := conn.NextReader()
		if err != nil {
			conn.Close()
			panic("Twig connection closed" + err.Error())

			// Todo: reconnect to twig
		}

		// Read payload
		var message GenericTwigMessage
		if err := conn.ReadJSON(&message); err != nil {
			println(utils.ParsedError(err, "Could not read message from twig", true))

			continue
		}

		switch message.Op {
		case 0:
			// Dispatch event (we don't need it :3)

		case 1:
			// Hola event handled above

		case 4:
			println("Heartbeat acknowledged", string(message.Data))

		default:
			println("unknown packet", string(message.Data))
		}
	}
}

func heartbeat(conn websocket.Conn, interval time.Duration) {
	for {
		if err := conn.WriteJSON(map[string]interface{}{
			"op": 3,
			"d": map[string]interface{}{
				"sent_at": time.Now().UnixMilli(),
			},
		}); err != nil {
			utils.ParsedError(err, "Could not start heartbeat to twig", true)
		}
		time.Sleep(interval)
	}
}
