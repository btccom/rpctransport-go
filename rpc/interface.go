package rpc

type Client interface {
	Request(req []byte) ([]byte, error)
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
