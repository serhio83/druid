package utils

import "bytes"

// ancient magic =)
// \u2388 Helm Symbol
// \u0fd5 Right-Facing Svasti Sign
// \u2693 Anchor
const decSymbol = "\u0fd5"
const space = ` `

//StringSplitter split \n & \r from byte slice & return string
func StringSplitter(b bytes.Buffer) string {
	var prepStr []byte
	for _, sm := range b.Bytes() {
		if string(sm) == "\n" || string(sm) == "\r" {
			prepStr = append(prepStr, space...)
		} else {
			prepStr = append(prepStr, sm)
		}
	}
	return string(prepStr)
}

//Envelope join string with prerfix & postfix decorators
func Envelope(s string) string {
	return decSymbol + space + s + space + decSymbol
}
