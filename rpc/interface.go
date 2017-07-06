package rpc

type Client interface {
	Request(req []byte) ([]byte, error)
	RequestAsync(req []byte) (<- chan []byte, <-chan error)
	Dial() error
	Close() error
}

type Server interface {
	Consume() <-chan ServerRequest
	Dial() error
	Close() error
}

type ServerRequest interface {
	Respond(response []byte) error
	Body() []byte
}
