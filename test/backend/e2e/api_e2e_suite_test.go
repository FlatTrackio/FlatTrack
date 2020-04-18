package api_e2e_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestApiE2e(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ApiE2e Suite")
}
