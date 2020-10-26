package uuid

import (
	"strings"

	"github.com/google/uuid"
)

type UUID = uuid.UUID

func New() UUID {
	return uuid.New()
}

func String() string {
	return strings.ReplaceAll(New().String(), "-", "")
}
