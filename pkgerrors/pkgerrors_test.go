// Copyright 2022 The go-analyzer Authors
// SPDX-License-Identifier: BSD-3-Clause

package pkgerrors_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/zchee/go-analyzer/pkgerrors"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, pkgerrors.Analyzer, "a", "b")
}
