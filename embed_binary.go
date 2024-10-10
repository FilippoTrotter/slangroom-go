//go:build generate
// +build generate

package slangroom

import _ "embed"

// Embedding Slangroom binary using go:embed
//
//go:embed slangroom-exec
var slangroomBinary []byte
