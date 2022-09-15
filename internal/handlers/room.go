package handlers

import (
	"fmt"
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

// TODO: Singleton

type Participant struct {
	Host bool
	Conn *websocket.Conn
	id   string
}

// roomsMap: [roomID string] -> []Participant
type roomsMap struct {
	mtx   sync.RWMutex
	rooms map[string]map[string]*Participant
}

// Initialize Rooms
func NewRoomsMap() *roomsMap {
	m := roomsMap{}
	m.rooms = make(map[string]map[string]*Participant)

	return &m
}

func (r *roomsMap) Get(id string) map[string]*Participant {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	return r.rooms[id]
}

func (r *roomsMap) CreateWithID(roomID string) string {
	r.mtx.Lock() // * Better lock writing style?

	r.rooms[roomID] = make(map[string]*Participant)

	r.mtx.Unlock()

	return roomID
}

func (r *roomsMap) Create() string {
	r.mtx.Lock() // * Better lock writing style?

	roomID := uuid.New().String()
	r.CreateWithID(roomID)

	return roomID
}

func (r *roomsMap) AddParticipant(roomID string, host bool, conn *websocket.Conn) string {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	p := Participant{host, conn, uuid.New().String()}

	r.rooms[roomID][p.id] = &p

	return p.id
}

func (r *roomsMap) RemoveParticipant(roomID string, pID string) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	delete(r.rooms[roomID], pID) // ? check if room exists?
}

func (r *roomsMap) Delete(roomID string) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	delete(r.rooms, roomID)
}

type BroadcastMsg struct {
	Body interface{} // TODO
	// Body *fiber.Map
	// Body   string
	RoomID string
	Client *websocket.Conn // ? Conn or Conn.Conn ?
}

var (
	_roomsMap = NewRoomsMap()
	broadcast = make(chan BroadcastMsg)
)

func broadcaster() {
	for {
		msg := <-broadcast
		fmt.Printf("Received a msg: %v\n", msg)

		for _, client := range _roomsMap.rooms[msg.RoomID] {
			// if client.Conn != msg.Client {
			err := client.Conn.WriteJSON(msg.Body)
			// err := client.Conn.WriteMessage(msg.Body)
			// err := client.Conn.WriteMessage(1, (msg.Body.([]uint8))) // * sends []byte
			if err != nil {
				log.Println(err)
				// client.Conn.Close() // ?
			}
			// }
		}
	}
}

func init() {
	fmt.Println("Init broadcaster")
	go broadcaster()
	// broadcaster() // ! fatal error: all goroutines are asleep - deadlock!
	fmt.Println("Passed broadcaster")
}

// TODO: Redirect ?
func CreateRoom(c *fiber.Ctx) error {
	roomID := _roomsMap.Create()

	return c.JSON(fiber.Map{"roomID": roomID})
}

// TODO: gofiber/websocket
// func JoinRoom(c *websocket.Conn) error {
// 	roomID := c.Params("id")
// 	if roomID == "" { // ? need to check ?
// 		return c.Status(fiber.ErrBadRequest.Code).SendString("Invalid roomID")
// 	}

// 	ws, err := c

// 	return nil
// }

func JoinRoom(c *websocket.Conn) /* void */ {
	roomID := c.Params("id")
	// if roomID == "" { // ? need to check ?
	// 	// return c.Status(fiber.ErrBadRequest.Code).SendString("Invalid roomID")
	// 	log.Println("No such room...")
	// 	return c.Close()
	// }

	if _, exists := _roomsMap.rooms[roomID]; !exists {
		log.Println("No such room...")
		// c.WriteMessage(fiber.ErrBadRequest.Code, []byte("No such room"))
		log.Printf("Available rooms: %v\n", _roomsMap.rooms)
		// c.Close()

		// return

		_roomsMap.CreateWithID(roomID)
		log.Printf("Entering room: %v", roomID)
	}

	pID := _roomsMap.AddParticipant(roomID, false, c) // ? host - false?

	for {
		// * Message is []byte which may be marshalled into struct / stringified JSON

		_, m, err := c.ReadMessage()
		// go func() {
		// 	log.Printf("Message: %v", m) // not working?
		// }()

		if err != nil {
			log.Printf("ReadMessage error: %v", err)
			_roomsMap.RemoveParticipant(roomID, pID) // ! Must remove Participants to not WriteMessage to nil Conn-s
			return                                   // * websocket: close 1001 (going away)
		}

		// broadcast <- BroadcastMsg{Body: &fiber.Map{"msg": m}, RoomID: roomID, Client: c}
		broadcast <- BroadcastMsg{Body: string(m), RoomID: roomID, Client: c}
		// broadcast <- BroadcastMsg{Body: m, RoomID: roomID, Client: c} // ! spits base64 with "/" if c.WriteJSON
	}

	// *
	// c.WriteJSON(fiber.Map{"MyMsg": "Oh hi Maark"})

	// ws, err := c

	// return nil
}
