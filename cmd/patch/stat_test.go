package main

import (
	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

var _ = Describe("Stat", func() {
	Context("when no patch", func() {
		It("returns error", func() {
			fs := afero.NewMemMapFs()
			source := "/source"
			err := afero.WriteFile(fs, source, []byte("source"), 0644)
			Expect(err).To(BeNil())

			var b bytes.Buffer
			err = StatInternal(fs, source, &b)
			Expect(err).ToNot(BeNil())

			b.String()
		})
	})
	Context("when patch", func() {})
})
