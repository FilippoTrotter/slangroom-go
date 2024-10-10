//go:build generate
// +build generate

package slangroom

import _ "embed"

// This will download the binary and store it as slangroom-exec
//
//go:generate sh -c "wget https://github.com/dyne/slangroom-exec/releases/latest/download/slangroom-exec-$(uname)-$(uname -m) -O ./slangroom-exec && chmod +x ./slangroom-exec"
