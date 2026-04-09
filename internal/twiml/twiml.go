package twiml

import (
	"encoding/xml"
	"fmt"
	"strings"
)

// EscapeForTwiML escapes text to be safe inside a TwiML XML element.
func EscapeForTwiML(s string) string {
	var b strings.Builder
	_ = xml.EscapeText(&b, []byte(s))
	return b.String()
}

// MessageResponse builds a minimal TwiML <Response><Message>...</Message></Response>.
func MessageResponse(text string) string {
	return fmt.Sprintf("<Response><Message>%s</Message></Response>", EscapeForTwiML(text))
}

