package bobber

type WebsocketService interface {
	HandleMessages()
	Broadcast(request *Request)
	AddClient(Client)
	DeleteClient(Client)
}

type Client interface {
	WriteMessage(int, []byte) error
	Close() error
}
