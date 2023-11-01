package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"lobby/server/handlers"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

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

func sendLobbyListen(p *ClientParams) error {
	log.Println("/lobby/listen")

	conn, res, err := websocket.DefaultDialer.Dial(
		p.webSocketUrlBase()+
			fmt.Sprintf("/lobby/listen/%s?player=%s", p.lobbyId, p.playerId),
		nil,
	)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusSwitchingProtocols {
		return errors.New("protocol switching does not work as expected")
	}
	defer res.Body.Close()
	defer conn.Close()

	return nil
}

func main() {
	params := &ClientParams{
		address:    "127.0.0.1",
		port:       "9990",
		playerName: "nekomimi",
		lobbyName:  "nekolobby",
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

	if err = sendLobbyCreate(params); err != nil {
		log.Panic(err)
	}

	if err = sendLobbyJoin(params); err != nil {
		log.Panic(err)
	}

	if err = sendLobbyListen(params); err != nil {
		log.Panic(err)
	}

	log.Println("done")
}
