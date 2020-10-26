// +build appengine

package rus

import (
	"io"
)

func checkIfTerminal(w io.Writer) bool {
	return true
}
