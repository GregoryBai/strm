const testMessage = document.querySelector('.test-message')
const messageContainer = document.querySelector('.messages')

const socket = new WebSocket('ws://localhost:8000/room/test')

const appendNewMessage = text => {
	messageContainer.innerHTML += `
            <div class="message">${text}</div>
`
}

socket.addEventListener('open', event => {
	socket.send('Hello Server!')
})

socket.addEventListener('close', event => {
	socket.send('Bye Server!')
})

// Listen for messages
socket.addEventListener('message', event => {
	// console.log('Message from server ', event.data)
	// console.log('🚀 ~ file: index.js ~ line 23 ~ data', event.data)
	// const data = JSON.parse(event.data)
	// console.log('🚀 ~ file: index.js ~ line 25 ~ data', data)
	// let message = atob(data.msg)
	// console.log('🚀 ~ file: index.js ~ line 27 ~ message', message)
	// appendNewMessage(message)
	appendNewMessage(event.data)
	// console.log({ response: event.data })
})

socket.addEventListener('error', error => {
	console.log('Socket error: ', error)
})

const input = document.querySelector('.input')
const button = document.querySelector('.button')

button.addEventListener('click', e => {
	const inputValue = input.value.trim()

	if (inputValue) {
		console.log('Sending: ', inputValue)
		socket.send(
			inputValue
			// JSON.stringify({
			// 	msg: inputValue,
			// })
		)
		input.value = ''
	}
})

// * ---------------------- WebRTC ---------------------- * //
// ? Should close ws connection upon completion ?

// * 1. Create PC instance with config for ICE servers
const pc = new RTCPeerConnection({
	// ? config not needed for local network ?
	iceServers: [
		{ urls: 'stun:stun.l.google.com:19302' },
		{
			urls: 'turn:openrelay.metered.ca:80',
			username: 'openrelayproject',
			credential: 'openrelayproject',
		},
		// {
		// 	urls: 'turn:openrelay.metered.ca:443',
		// 	username: 'openrelayproject',
		// 	credential: 'openrelayproject',
		// },
	],
})

// * 2. Prepare data channel (? negotiationneeded event) with its even listeners
// ** 2.1 Might want to use
const dataChannel = pc.createDataChannel('default')

// **** RTC WebSocket (as a way to exchange SDPs between parties)

const RTCSocket = new WebSocket('ws://localhost:8000/rtc/ws')

dataChannel.addEventListener('error', e => console.log('Error: ', e)) // ? useless ?
dataChannel.addEventListener('open', e => console.log('Started: ', e))
dataChannel.addEventListener('message', e => console.log('Message: ', e))
dataChannel.addEventListener('close', e => console.log('Close: ', e))

RTCSocket.addEventListener('error', e => console.log('RTCSocket error: ', e, e.error))
RTCSocket.addEventListener('open', e => {
	console.log('RTCSocket opened: ', e)
})
RTCSocket.addEventListener('close', e => console.log('RTCSocket closed: ', e))
RTCSocket.addEventListener('message', e => {
	const message = JSON.parse(e.data)

	if (message.type === 'answer') {
		console.log('Message (answer): ', message)
		pc.setRemoteDescription(message.data)

		// dataChannel.send('Woooow') // ?
	} else {
		console.log('Message (not answer): ', message)
	}
})

// * 3.1. Gather ICE candidates that update SDP
// pc.addEventListener('icecandidate', e => console.log('ICE Candidate: ', e))
pc.addEventListener('icecandidateerror', e => console.log('ICE Candidate Error: ', e))

// * 3.2. Wait
setTimeout(() => {
	//

	// * 4. Create offer to be sent to the server & set it as LD
	pc.createOffer()
		.then(desc => {
			pc.setLocalDescription(desc)
			// ! Might throw error on accessing unavailable desc
			RTCSocket.send(JSON.stringify({ type: 'offer', data: desc }))
		})
		.finally(() => console.log('Offer Created'))
		.catch(console.error)
}, 1000)

// * 5. Receive answer from the other party and set it as RD

// ✅✅✅

pc.addEventListener('datachannel', e => {
	console.log('Data Channel: ', e)
})
/*
* Other party needs to:
	
	* - NewRTCPeerConnection ( defer pc.close() & watch for connectionstatechange)
	* - Gather ICE candidates
	* - Set up Data Channel initiated by a party 'ondatachannel' & add message listeners to it
	* - Set Remote Description with the offer received
	* - Create an answer to send to the party and set it as Local Description

* Data channels are available when both partie's LDs & RDs have been set
*/
