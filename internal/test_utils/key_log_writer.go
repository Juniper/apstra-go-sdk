// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package testutils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const keyLogFile = "SSLKEYLOGFILE"

// KeyLogWriterFromEnv takes an environment variable which might name a logfile for
// exporting TLS session keys. If so, it returns an io.Writer to be used for
// that purpose, and the name of the logfile file.
func KeyLogWriterFromEnv(t testing.TB) *os.File {
	t.Helper()
	
	fileName, foundKeyLogFile := os.LookupEnv(keyLogFile)
	if !foundKeyLogFile {
		return nil
	}

	// expand ~ style home directory
	if strings.HasPrefix(fileName, "~/") {
		dirname, _ := os.UserHomeDir()
		fileName = filepath.Join(dirname, fileName[2:])
	}

	err := os.MkdirAll(filepath.Dir(fileName), os.FileMode(0o600))
	if err != nil {
		t.Fatalf("Error creating keylog dir: %v", err)
	}

	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		t.Fatalf("Error opening keylog file: %v", err)
	}

	return f
}
