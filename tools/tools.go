//go:build tools

package tools

import (
	// license compliance
	_ "github.com/chrismarget-j/go-licenses"

	//opinionated code formatting
	_ "mvdan.cc/gofumpt"
)
