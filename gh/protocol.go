package gh

// Constents
const (
	IPHeaderLength  int = 20
	UDPHeaderLength int = 8
	GHHeaderLength  int = 4

	GHFlagHeart    uint8 = 0x01 // Heart beat
	GHFlagRouteReq uint8 = 0x02 // Route request
	GHFlagAuth     uint8 = 0x04 // Auth request
	GHFlagData     uint8 = 0x08 // Data package
	GHFlagClose    uint8 = 0x10 // Close connection request
)

type groundhogHeader struct {
	flag   uint8
	index  uint8
	length uint16
}

type groundhogPackage struct {
	groundhogHeader
	payload []byte
}
