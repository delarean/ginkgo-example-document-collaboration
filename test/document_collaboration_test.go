package collaboration_tests

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCollaborationService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Collaboration Service Suite")
}
