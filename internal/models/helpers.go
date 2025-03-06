package models

import "github.com/lithammer/shortuuid"

func GenerateID() string {
	return shortuuid.New()[:12]
}

func GenerateIDWithLength(length int) string {
	return shortuuid.New()[:length]
}
