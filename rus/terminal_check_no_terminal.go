// +build js nacl plan9

package rus

import (
	"io"
)

func checkIfTerminal(w io.Writer) bool {
	return false
}
