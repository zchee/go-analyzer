// Copyright 2022 The go-analyzer Authors
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"github.com/zchee/go-analyzer/pkgerrors"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(pkgerrors.Analyzer) }
