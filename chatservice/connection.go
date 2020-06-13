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
		Action: "message",
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

func (cc *ChatConn) Join(s *server, room string) (*nats.Subscription, error) {
	s.CreateRoomIfNotExist(room)

	token, err := GetToken(map[string]interface{}{
		"id": cc.ID,
	})
	if err != nil {
		log.WithError(err).Error("could not create token")
	}

	sub, err := s.nc.Subscribe(room, func(msg *nats.Msg) {
		var reqMsg Req
		err := json.Unmarshal(msg.Data, &reqMsg)
		if err != nil {
			log.WithError(err).Error("could not unmarshal message from queue")
		}

		if reqMsg.ID == cc.ID {
			return
		}

		m, err := json.Marshal(Resp{
			Token: token,
			ID: reqMsg.ID,
			Message: reqMsg.Message,
			Action: "message",
		})
		if err != nil {
			log.WithError(err).Error("could not marshal message to send")
		}

		err = cc.Conn.WriteMessage(websocket.TextMessage, m)
		if err != nil {
			log.WithError(err).Error("could not write to websocket")
		}
	})

	if err != nil {
		return nil, fmt.Errorf("subscribe:%w", err)
	}

	m, err := json.Marshal(Resp{
		Token: token,
		ID: cc.ID,
		Action: "joined",
	})
	if err != nil {
		return sub, fmt.Errorf("could not marshal:%w", err)
	}

	err = cc.Conn.WriteMessage(websocket.TextMessage, m)
	if err != nil {
		return sub, fmt.Errorf("write subscribe:%w", err)
	}

	return sub, nil
}
