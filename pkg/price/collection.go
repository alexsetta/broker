package price

const (
	limit = 100
)

type Collection struct {
	id     string
	prices []float64
}

// NewCollection returns a new Collection struct
func NewCollection(id string) *Collection {
	return &Collection{
		id:     id,
		prices: []float64{},
	}
}

// Add appends a new price to the prices slice
func (c *Collection) Add(price float64) {
	if len(c.prices) == limit {
		c.prices = c.prices[1:]
	}
	c.prices = append(c.prices, price)
}
