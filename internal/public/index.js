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
	console.log('Message from server ', event.data)
	// console.log('🚀 ~ file: index.js ~ line 23 ~ data', event.data)
	// const data = JSON.parse(event.data)
	// console.log('🚀 ~ file: index.js ~ line 25 ~ data', data)
	// let message = atob(data.msg)
	// console.log('🚀 ~ file: index.js ~ line 27 ~ message', message)
	// appendNewMessage(message)
	appendNewMessage(event.data)
	console.log({ response: event.data })
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
