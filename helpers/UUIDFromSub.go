package helpers

import (
	"github.com/google/uuid"
)

func UUIDFromSub(sub string) uuid.UUID {
	// Use a fixed namespace UUID (e.g., UUID for DNS)
	namespace := uuid.NameSpaceDNS
	// Generate UUIDv5 based on namespace + sub string
	return uuid.NewSHA1(namespace, []byte(sub))
}
