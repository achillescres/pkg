package oid

import (
	"strconv"
)

type ID string

const Undefined ID = ""

// ToID casts string to id
func ToID(id string) ID {
	return ID(id)
}

// AssertID claim any type and determine whether it satisfies oid.ID type contract
// if so returns (ID, true) else returns ("", false)
func AssertID(id any) (ID, bool) {
	nID, ok := id.(string)
	if !ok {
		return "", false
	}

	if IsID(nID) {
		return ID(nID), true
	} else {
		return "", false
	}
}

func IsID(id string) bool {
	if IsUndefined(ID(id)) {
		return false
	}
	if idNum, err := strconv.Atoi(id); err != nil || idNum < 1 {
		return false
	}
	return true
}

func IsUndefined(id ID) bool {
	return id == Undefined
}

func (id ID) String() string {
	return string(id)
}
