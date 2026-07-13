package provider

import "github.com/google/uuid"

type IDGenerator interface {
	Generate() string
}

type uuidGenerator struct{}

func NewUUIDGenerator() IDGenerator {
	return &uuidGenerator{}
}

func (g *uuidGenerator) Generate() string {
	return uuid.New().String()
}
