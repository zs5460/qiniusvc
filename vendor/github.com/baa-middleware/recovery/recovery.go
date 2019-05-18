// Package recovery provider a baa middleware which recovers from panics anywhere in the chain.
package recovery

import (
	"fmt"
	"runtime"

	"gopkg.in/baa.v1"
)

// Recovery returns a baa middleware which recovers from panics anywhere in the chain
// and handles the control to the centralized HTTPErrorHandler.
func Recovery() baa.HandlerFunc {
	return func(c *baa.Context) {
		defer func() {
			if err := recover(); err != nil {
				trace := make([]byte, 1<<16)
				n := runtime.Stack(trace, true)
				c.Error(fmt.Errorf("panic recover\n %v\n stack trace %d bytes\n %s", err, n, trace[:n]))
			}
		}()

		c.Next()
	}
}
