package oid

import (
	"strconv"
	"strings"
)

type ID string

const Undefined ID = "-1"

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
	if intID, err := strconv.Atoi(id); err == nil && !IsUndefined(id) && intID > 0 {
		return true
	}
	return false
}

func IsUndefined(id ID) bool {
	return id == Undefined || strings.TrimSpace(string(id)) == ""
}

func (id ID) String() string {
	return string(id)
}
