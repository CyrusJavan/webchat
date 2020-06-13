package chatservice

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
)

type KeyValueStore interface {
	Get(string) (string, error)
	Set(string, string) error
}

type server struct {
	kv KeyValueStore
	nc *nats.Conn
	r *mux.Router
}

type Req struct {
	Action string
	Room string
	Message string
	ID string
	Token string
}

type Resp struct {
	ID string `json:"id,omitempty"`
	Token string `json:"token,omitempty"`
	Message string `json:"message,omitempty"`
	Action string `json:"action,omitempty"`
}

func NewServer(nc *nats.Conn) *server {
	s := &server{
		kv: NewMapStore(),
		nc: nc,
		r: mux.NewRouter(),
	}

	s.r.HandleFunc("/", HomeHandler)
	s.r.Path("/chat").
		HandlerFunc(s.ChatHandler)

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.r.ServeHTTP(w, r)
}

const KB = 1024
var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1*KB,
	WriteBufferSize: 1*KB,
}

func (s *server) ChatHandler(w http.ResponseWriter, r *http.Request) {
	c, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	var sub *nats.Subscription
	defer c.Close()
	closer := c.CloseHandler()
	c.SetCloseHandler(func(code int, text string) error {
		err := sub.Drain()
		if err != nil {
			log.WithError(err).Error("draining subscription issue")
		}

		return closer(code, text)
	})

	cc := &ChatConn{
		ID: uuid.New().String(),
		Conn: c,
	}

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("recv: %s", message)

		var req Req
		err = json.Unmarshal(message, &req)
		if err != nil {
			log.Println("json unmarshal:", err)
			return
		}

		switch req.Action {
		case "join":
			sub, err = cc.Join(s, req.Room)
			if err != nil {
				log.WithError(err).Error("could not join")
				if sub != nil {
					if err := sub.Unsubscribe(); err != nil {
						return
					}
				}
			}
		case "send":
			claims, err := GetClaims(req.Token)
			if err != nil {
				return
			}
			claimedID := claims["id"].(string)

			if claimedID != req.ID {
				log.Errorf("client sent wrong id. claimed:%s actual:%s", claimedID, req.ID)
				return
			}

			if err := cc.Send(s, req.Room, req.Message); err != nil {
				log.WithError(err).Error("could not send")
				return
			}
		}
	}
}

var homeTemplate = template.Must(template.New("home.tpl").ParseFiles("./tpl/home.tpl"))
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if err := homeTemplate.Execute(w, nil); err != nil {
		log.WithError(err).Fatal("could not execute home template")
	}
}

func (s *server) CreateRoomIfNotExist(name string) {
	_, err := s.kv.Get(name)
	if err == nil {
		return
	}
	_ = s.kv.Set(name, name)
}
