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
