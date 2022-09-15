package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/websocket/v2"
	"github.com/pion/webrtc/v3"
)

const (
	EventTypeOffer  = "offer"
	EventTypeAnswer = "answer"
)

type Event struct {
	Type string                    `json:"type"`
	Data webrtc.SessionDescription `json:"data"`
}

var (
	logCh = make(chan interface{}, 2)

	pc *webrtc.PeerConnection
)

func init() {
	go func() {
		for {
			fmt.Printf("Event: %v", <-logCh)
		}
	}()

	config := webrtc.Configuration{

		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{
					"stun:stun.l.google.com:19302",
				},
			},
			{
				URLs: []string{
					"turn:openrelay.metered.ca:80",
					// "turn:openrelay.metered.ca:443",
				},
				Username:   "openrelayproject",
				Credential: "openrelayproject",
				// CredentialType: webrtc.ICECredentialTypePassword,
			},
		},
	}

	pc_, err := webrtc.NewPeerConnection(config)
	if err != nil {
		fmt.Println(err)
	}
	pc = pc_

	pc.OnDataChannel(func(dc *webrtc.DataChannel) {
		fmt.Printf("%v+", dc)
	})

	// defer func() {
	// 	if err := pc.Close(); err != nil {
	// 		fmt.Printf("Can't close connection %v\n", err)
	// 	}
	// }()
}

func InitWebRTC(c *websocket.Conn) {
	defer c.Close()

	for {
		var e Event
		if err := c.ReadJSON(&e); err != nil {
			fmt.Printf("InitWebRTC error: %v\n", err)
			// continue
			break
		}

		switch e.Type {
		case EventTypeOffer:
			pc.SetRemoteDescription(e.Data)

			// * gatherComplete is a channel that will be closed when
			// * the gathering of local candidates is complete.
			// gatherComplete := webrtc.GatheringCompletePromise(pc)
			// <-gatherComplete

			<-time.After(1 * time.Second)

			answer, err := pc.CreateAnswer(nil)
			if err != nil {
				fmt.Printf("pc.CreateAnswer err: %v\n", err)
			}

			pc.SetLocalDescription(answer)

			err = c.WriteJSON(Event{Type: EventTypeAnswer, Data: answer})
			if err != nil {
				fmt.Printf("c.WriteJSON err: %v\n", err)
				continue
			}

			// dc, err := pc.CreateDataChannel("default", nil)
			// if err != nil {
			// 	fmt.Printf("pc.CreateDataChannel err: %v\n", err)
			// }

			// <-time.After(1 * time.Second) //

			// if err := dc.SendText("Huey"); err != nil {
			// 	fmt.Printf("dc.SendText err: %v\n", err)
			// }

			fmt.Printf("Offer Created!\n")
		default:
			fmt.Printf("Event: %v\n", e)
		}
	}
}
