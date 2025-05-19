package cred

import "time"

type Cred[T any] struct {
	Cred      T
	ExpiresAt time.Time
	Ttl       time.Duration
}

func (c *Cred[T]) Expired() bool {
	return !c.ExpiresAt.After(time.Now())
}

func (c *Cred[T]) Update(new T) {

}
