// Copyright 2018 The Rind Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
)

// packet to string
func pString(p Packet) string {
	return fmt.Sprint(p.message.ID)
}
