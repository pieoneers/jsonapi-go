// Copyright (c) 2020 Pieoneers Software Incorporated. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package jsonapi_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestJSONAPI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "JSONAPI Suite")
}
