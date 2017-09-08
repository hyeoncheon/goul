package pipes

// ItemGeneric is a structure for the generic byte slice data.
type ItemGeneric struct {
	meta string
	data []byte
}

// String implements goul.Item
func (c *ItemGeneric) String() string {
	return c.meta
}

// Data implements goul.Item
func (c *ItemGeneric) Data() []byte {
	return c.data
}
