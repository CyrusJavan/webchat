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
}

type Resp struct {
	ID string `json:"id"`
	Message string `json:"message"`
	Subscribed bool `json:"subscribed"`
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
	defer c.Close()

	cc := &ChatConn{
		ID: uuid.New().String(),
		Conn: c,
	}

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		var req Req
		err = json.Unmarshal(message, &req)
		if err != nil {
			log.Println("json unmarshal:", err)
		}

		switch req.Action {
		case "join":
			if err := cc.Join(s, req.Room); err != nil {
				log.WithError(err).Error("could not join")
			}
		case "send":
			if err := cc.Send(s, req.Room, req.Message); err != nil {
				log.WithError(err).Error("could not send")
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
