package bytes

func ToUint16(h byte,l byte) uint16 {
	return  uint16(h)<<8 | uint16(l)

}