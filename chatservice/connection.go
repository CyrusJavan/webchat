package chatservice

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)

type ChatConn struct {
	ID string
	Conn *websocket.Conn
}

func (cc *ChatConn) Send(s *server, room, message string) error {
	s.CreateRoomIfNotExist(room)

	m, err := json.Marshal(Resp{
		ID: cc.ID,
		Message: message,
	})
	if err != nil {
		return fmt.Errorf("could not marshal:%w", err)
	}
	log.Printf("sent: %s", m)

	err = s.nc.Publish(room, m)
	if err != nil {
		return fmt.Errorf("publish:%w", err)
	}

	return nil
}

func (cc *ChatConn) Join(s *server, room string) error {
	s.CreateRoomIfNotExist(room)

	_, err := s.nc.Subscribe(room, func(msg *nats.Msg) {
		var reqMsg Req
		err := json.Unmarshal(msg.Data, &reqMsg)
		if err != nil {
			log.Println("json unmarshal:", err)
		}

		if reqMsg.ID == cc.ID {
			return
		}

		m, err := json.Marshal(Resp{
			ID: reqMsg.ID,
			Message: reqMsg.Message,
			Subscribed: false,
		})
		if err != nil {
			log.WithError(err).Print("could not marshal")
		}

		err = cc.Conn.WriteMessage(websocket.TextMessage, m)
		if err != nil {
			log.Println("write:", err)
		}

		log.Println(reqMsg)
	})

	if err != nil {
		return fmt.Errorf("subscribe:%w", err)
	}

	m, err := json.Marshal(Resp{
		ID: cc.ID,
		Subscribed: true,
	})
	if err != nil {
		return fmt.Errorf("could not marshal:%w", err)
	}

	err = cc.Conn.WriteMessage(websocket.TextMessage, m)
	if err != nil {
		return fmt.Errorf("write subscribe:%w", err)
	}

	return nil
}
