// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package pointer_test

import (
	"net"
	"testing"

	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	"github.com/stretchr/testify/require"
)

func TestZeroOf(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		d := ""
		p := pointer.To(d)
		r := pointer.ZeroOf(p)
		require.IsType(t, d, r)
		require.Zero(t, r)
	})

	t.Run("struct", func(t *testing.T) {
		d := net.IPAddr{}
		p := pointer.To(d)
		r := pointer.ZeroOf(p)
		require.IsType(t, d, r)
		require.Zero(t, r)
	})

	t.Run("nil_pointer", func(t *testing.T) {
		p := (*net.IPAddr)(nil)
		r := pointer.ZeroOf(p)
		require.IsType(t, net.IPAddr{}, r)
		require.Zero(t, r)
	})
}

func TestTypeOf(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		d := ""
		p := pointer.To(d)
		r := pointer.TypeOf(p)
		require.Equal(t, "string", r)
	})

	t.Run("struct", func(t *testing.T) {
		d := net.IPAddr{}
		p := pointer.To(d)
		r := pointer.TypeOf(p)
		require.Equal(t, "net.IPAddr", r)
	})

	t.Run("nil_pointer", func(t *testing.T) {
		p := (*net.IPAddr)(nil)
		r := pointer.TypeOf(p)
		require.Equal(t, "net.IPAddr", r)
	})
}
