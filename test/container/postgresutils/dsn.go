package postgresutils

import "fmt"

func DSN(host string, port uint, db, user, password string) string {
	return fmt.Sprintf("postgresql://%s:%d/%s?user=%s&password=%s", host, port, db, user, password)
}
