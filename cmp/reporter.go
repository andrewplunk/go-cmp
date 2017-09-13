// Copyright 2017, The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package cmp

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/andrewplunk/go-cmp/cmp/internal/value"
)

type defaultReporter struct {
	Option
	diffs  []string // List of differences, possibly truncated
	ndiffs int      // Total number of differences
	nbytes int      // Number of bytes in diffs
	nlines int      // Number of lines in diffs
}

var _ reporter = (*defaultReporter)(nil)

func (r *defaultReporter) Report(x, y reflect.Value, eq bool, p Path) {
	if eq {
		// TODO: Maybe print some equal results for context?
		return // Ignore equal results
	}
	const maxBytes = 4096
	const maxLines = 256
	r.ndiffs++
	if r.nbytes < maxBytes && r.nlines < maxLines {
		sx := value.Format(x, true)
		sy := value.Format(y, true)
		if sx == sy {
			// Stringer is not helpful, so rely on more exact formatting.
			sx = value.Format(x, false)
			sy = value.Format(y, false)
		}

		// Add type information if available.
		var tx, ty string
		if x.IsValid() {
			tx = fmt.Sprintf("%T:", x.Interface())
		}

		if y.IsValid() {
			ty = fmt.Sprintf("%T:", y.Interface())
		}

		s := fmt.Sprintf("%#v:\n\t-: %s%s\n\t+: %s%s\n", p, tx, sx, ty, sy)
		r.diffs = append(r.diffs, s)
		r.nbytes += len(s)
		r.nlines += strings.Count(s, "\n")
	}
}

func (r *defaultReporter) String() string {
	s := strings.Join(r.diffs, "")
	if r.ndiffs == len(r.diffs) {
		return s
	}
	return fmt.Sprintf("%s... %d more differences ...", s, len(r.diffs)-r.ndiffs)
}
