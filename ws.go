package bobber

type WebsocketService interface {
	HandleMessages()
	Broadcast() chan *RequestDetail
	AddClient(Client)
	DeleteClient(Client)
}

type Client interface {
	WriteMessage(int, []byte) error
	Close() error
}
