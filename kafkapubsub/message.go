package kafkapubsub

// Message is interface for kafka topics use
type Message interface {
	// Bytes converts Message to []byte
	Bytes() ([]byte, error)

	// Unmarshal scans []byte and unmarshals and returns it
	Unmarshal([]byte) (Message, error)

	// Key returns Message's key
	Key() string
}
