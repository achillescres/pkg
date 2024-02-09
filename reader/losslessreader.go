package reader

import "io"

func ReadLossless(r io.Reader) ([]byte, error) {
	data := make([]byte, 0)
	buffer := make([]byte, 1024)
	for {
		n, err := r.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n == 0 {
			break
		}
		data = append(data, buffer[:n]...)
	}
	return data, nil
}
