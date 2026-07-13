package provider

import (
	"strings"

	"github.com/google/uuid"
)

type IDGenerator interface {
	Generate() string
	GeneratePublic() string
}

type idGenerator struct{}

func NewIdGenerator() IDGenerator {
	return &idGenerator{}
}

func (g *idGenerator) Generate() string {
	return uuid.New().String()
}

func (g *idGenerator) GeneratePublic() string {
	s := strings.ToUpper(strings.ReplaceAll(uuid.New().String(), "-", ""))

	return s[0:3] + "-" + s[3:7] + "-" + s[7:11] + "-" + s[11:15] + "-" + s[15:19]
}
