package kafkapubsub

// Message is interface for kafka topics use
type Message interface {
	// Bytes converts Message to []byte
	Bytes() []byte
	// Scan scans []byte and unmarshals it to Message
	Scan([]byte) error
}
