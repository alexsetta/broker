package price

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestCollection_Add(t *testing.T) {
	c := NewCollection("teste")
	assert.NotNil(t, c, "The collection should not be nil")

	rand.Seed(0)
	for i := 0; i < 101; i++ {
		c.Add(995 + rand.Float64()*10)
	}

	assert.Equal(t, 100, len(c.prices), "The length of the prices slice should be 100")
}
