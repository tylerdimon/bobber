package bobber

type WebsocketService interface {
	HandleMessages()
	Broadcast() chan string
	AddClient(Client)
	DeleteClient(Client)
}

type Client interface {
	WriteMessage(int, []byte) error
	Close() error
}
