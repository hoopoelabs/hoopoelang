// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"github.com/hoopoelabs/hoopoelang"
)

type LogBeginCallback func(testName string, testFile *TestFile)
type LogProblemsCallback func(testName string, testFile *TestFile, diags hcl.Diagnostics)
