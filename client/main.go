package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"lobby/lobby/message"
	"lobby/server/handlers"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
)

type MessageListener struct {
	conn *websocket.Conn

	msgCh chan string
	errCh chan error
}

func NewMessageListener(conn *websocket.Conn) *MessageListener {
	return &MessageListener{
		conn:  conn,
		msgCh: make(chan string, 10),
		errCh: make(chan error),
	}
}

func (l *MessageListener) listen() {
	format := "%s: %s"
	defer l.Close()
LOOP:
	for {
		envelope := message.Envelope{}
		err := l.conn.ReadJSON(&envelope)
		if err != nil {
			l.errCh <- err
			break LOOP
		}

		for _, msg := range envelope.Messages {
			l.msgCh <- fmt.Sprintf(format, msg.Key, msg.Value)
		}
	}
}

func (l *MessageListener) Close() error {
	return l.conn.Close()
}

type ClientParams struct {
	address string
	port    string

	playerName string
	playerId   string
	lobbyName  string
	lobbyId    string
}

func (p *ClientParams) httpUrlBase() string {
	return fmt.Sprintf("http://%s:%s", p.address, p.port)
}

func (p *ClientParams) webSocketUrlBase() string {
	return fmt.Sprintf("ws://%s:%s", p.address, p.port)
}

func sendLobbyCreate(p *ClientParams) error {
	log.Println("/lobby/create")

	form := url.Values{
		"lobby-name": {p.lobbyName},
	}
	httpRes, err := http.PostForm(
		p.httpUrlBase()+"/lobby/create",
		form,
	)
	if err != nil {
		return err
	}
	defer httpRes.Body.Close()

	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return err
	}
	log.Println(string(body))

	res := &handlers.LobbyCreateResponse{}
	if err = json.Unmarshal(body, res); err != nil {
		return err
	}

	p.lobbyId = res.LobbyId
	return nil
}

func sendLobbyJoin(p *ClientParams) error {
	log.Println("/lobby/join")

	form := url.Values{
		"lobby-id":    {p.lobbyId},
		"player-name": {p.playerName},
	}
	httpRes, err := http.PostForm(
		p.httpUrlBase()+"/lobby/join",
		form,
	)
	if err != nil {
		return err
	}
	defer httpRes.Body.Close()

	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return err
	}
	log.Println(string(body))

	res := &handlers.LobbyJoinResponse{}
	if err = json.Unmarshal(body, res); err != nil {
		return err
	}

	p.playerId = res.PlayerId
	return nil
}

func sendLobbyListen(p *ClientParams) (*MessageListener, error) {
	log.Println("/lobby/listen")

	conn, res, err := websocket.DefaultDialer.Dial(
		p.webSocketUrlBase()+
			fmt.Sprintf("/lobby/listen/%s?player=%s", p.lobbyId, p.playerId),
		nil,
	)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusSwitchingProtocols {
		return nil, errors.New("protocol switching does not work as expected")
	}
	defer res.Body.Close()

	return NewMessageListener(conn), nil
}

func main() {
	if len(os.Args) <= 1 {
		log.Fatal("use\n-c for create new lobby\n-l ID for join existing lobby")
	}

	params := &ClientParams{
		address:    "127.0.0.1",
		port:       "9990",
		playerName: "nekomimi",
		lobbyName:  "nekolobby",
	}

	doCreate := os.Args[1] == "-c"
	if os.Args[1] == "-l" {
		params.lobbyId = os.Args[2]
	}

	log.Printf(
		"start sending requests to %s:%s",
		params.address,
		params.port,
	)

	log.Println("/")
	httpRes, err := http.Get(params.httpUrlBase())
	if err != nil {
		log.Panic(err)
	}
	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		log.Panic(err)
	}
	log.Println(string(body))
	httpRes.Body.Close()

	if doCreate {
		if err = sendLobbyCreate(params); err != nil {
			log.Panic(err)
		}
	}

	if err = sendLobbyJoin(params); err != nil {
		log.Panic(err)
	}

	listener, err := sendLobbyListen(params)
	if err != nil {
		log.Panic(err)
	}

	go listener.listen()

	for {
		select {
		case msg := <-listener.msgCh:
			log.Println(msg)
		case err := <-listener.errCh:
			log.Panic(err)
		}
	}
}
