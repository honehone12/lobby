package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"lobby/lobby/message"
	"lobby/server/handlers"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type Connection struct {
	conn *websocket.Conn

	msgCh chan string
	errCh chan error
}

func NewMessageListener(conn *websocket.Conn) *Connection {
	return &Connection{
		conn:  conn,
		msgCh: make(chan string, 10),
		errCh: make(chan error),
	}
}

func (c *Connection) listen() {
	format := "%s: %s"
	defer c.Close()
LOOP:
	for {
		envelope := message.Envelope{}
		err := c.conn.ReadJSON(&envelope)
		if err != nil {
			c.errCh <- err
			break LOOP
		}

		for _, msg := range envelope.Messages {
			c.msgCh <- fmt.Sprintf(format, msg.Key, msg.Value)
		}
	}
}

func (c *Connection) SendChatMessage() {
	ticker := time.Tick(time.Second)
LOOP:
	for range ticker {
		envelope := message.NewChatMessageRequest("nyannyan")
		err := c.conn.WriteJSON(envelope)
		if err != nil {
			c.errCh <- err
			break LOOP
		}
	}
}

func (c *Connection) Close() error {
	return c.conn.Close()
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

func sendLobbyListen(p *ClientParams) (*Connection, error) {
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
	createFlag := flag.Bool("c", false, "if true create new lobby")
	sendChat := flag.Bool("s", false, "if true send chat message")
	lobbyId := flag.String("j", "", "lobby id for join")
	lobbyNama := flag.String("l", "nekolobby", "lobby name for create")
	playerName := flag.String("n", "nekomimi", "player name")

	flag.Parse()

	params := &ClientParams{
		address:    "127.0.0.1",
		port:       "9990",
		playerName: *playerName,
		playerId:   "",
		lobbyName:  *lobbyNama,
		lobbyId:    *lobbyId,
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

	if *createFlag {
		if err = sendLobbyCreate(params); err != nil {
			log.Panic(err)
		}
	}

	if err = sendLobbyJoin(params); err != nil {
		log.Panic(err)
	}

	conn, err := sendLobbyListen(params)
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	go conn.listen()

	if *sendChat {
		go conn.SendChatMessage()
	}

	for {
		select {
		case msg := <-conn.msgCh:
			log.Println(msg)
		case err := <-conn.errCh:
			log.Panic(err)
		}
	}
}
