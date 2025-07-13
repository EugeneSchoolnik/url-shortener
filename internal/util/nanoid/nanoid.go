package nanoid

import gonanoid "github.com/matoous/go-nanoid"

type IDGenerator struct {
	alphabet string
	size     int
}

func New(alphabet string, size int) *IDGenerator {
	return &IDGenerator{alphabet, size}
}

func (g *IDGenerator) ID() (string, error) {
	return gonanoid.Generate(g.alphabet, g.size)
}
