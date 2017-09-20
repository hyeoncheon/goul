package goul

// constants...
const (
	ItemTypeUnknown   = "unknown"
	ItemTypeRawPacket = "rawpacket"
)

//** types for goul, items ------------------------------------------

// Item is an interface for passing data between pipes.
type Item interface {
	String() string
	Data() []byte
}

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
