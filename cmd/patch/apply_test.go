package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/shim/pkg/patch"
	"github.com/spf13/afero"
)

var _ = Describe("Apply", func() {
	Context("when no patch", func() {
		It("creates the patch", func() {
			fs := afero.NewMemMapFs()
			source := "/source"
			target := "/target"
			data := "/data"

			err := afero.WriteFile(fs, source, []byte("source"), 0644)
			Expect(err).To(BeNil())

			dataContent := "data"
			err = afero.WriteFile(fs, data, []byte(dataContent), 0644)
			Expect(err).To(BeNil())

			err = ApplyInternal(fs, source, target, data)
			Expect(err).To(BeNil())

			ok, err := afero.Exists(fs, target)
			Expect(err).To(BeNil())
			Expect(ok).To(BeTrue())

			targetFile, err := fs.Open(target)
			Expect(err).To(BeNil())
			defer targetFile.Close()

			p, err := patch.Get(targetFile)
			Expect(err).To(BeNil())
			Expect(p).ToNot(BeNil())
			Expect(p.Data).To(Equal([]byte(dataContent)))
		})
	})
	Context("when existing patch", func() {})
})
