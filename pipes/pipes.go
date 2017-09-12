package pipes

// ItemGeneric is a structure for the generic byte slice data.
type ItemGeneric struct {
	Meta string
	DATA []byte
}

// String implements goul.Item
func (c *ItemGeneric) String() string {
	return c.Meta
}

// Data implements goul.Item
func (c *ItemGeneric) Data() []byte {
	return c.DATA
}
