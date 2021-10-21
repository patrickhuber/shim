package patch_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/shim/pkg/patch"
	"github.com/spf13/afero"
)

var _ = Describe("Patch", func() {
	Describe("FromString", func() {
		It("returns the correct sizes", func() {
			data := "{ \"name\" : \"test\" }"
			p := patch.FromString(data)
			Expect(p).ToNot(BeNil())
			Expect(len(p.Data)).To(Equal(len(data)))
			Expect(int(p.Size)).To(Equal(len(data)))
		})
	})
	Describe("Apply", func() {
		Context("when empty file", func() {
			It("is applied", func() {

				content := "this is some content"
				fs := afero.NewMemMapFs()
				err := afero.WriteFile(fs, "/target", []byte(content), 0644)
				Expect(err).To(BeNil())

				data := "[]"
				p := patch.FromString(data)
				Expect(p).ToNot(BeNil())

				target, err := fs.OpenFile("/target", os.O_APPEND|os.O_WRONLY, 0644)
				Expect(err).To(BeNil())

				err = p.Apply(target)
				Expect(err).To(BeNil())
				defer target.Close()

				newPatch, err := patch.Get(target)
				Expect(err).To(BeNil())
				Expect(newPatch).ToNot(BeNil())
				Expect(newPatch.Size).To(Equal(uint32(len(newPatch.Data))))
				Expect(string(newPatch.Data)).To(Equal(data))
			})
		})
	})
	Describe("Get", func() {
		Context("when no patch", func() {
			It("returns null", func() {
				fs := afero.NewMemMapFs()
				err := afero.WriteFile(fs, "/source", []byte{}, 0644)
				Expect(err).To(BeNil())

				f, err := fs.Open("/source")
				Expect(err).To(BeNil())

				p, err := patch.Get(f)
				Expect(err).To(BeNil())
				Expect(p).To(BeNil())
			})
		})
		Context("when patch", func() {
			It("returns patch", func() {

				// data(data size) | data size(4) | magic(2) |
				bytes := []byte{0x67, 0x6f, 0x00, 0x00, 0x00, 0x02, 0x21, 0x23}

				fs := afero.NewMemMapFs()
				err := afero.WriteFile(fs, "/source", bytes, 0644)
				Expect(err).To(BeNil())

				source, err := fs.Open("/source")
				Expect(err).To(BeNil())
				defer source.Close()

				p, err := patch.Get(source)
				Expect(err).To(BeNil())
				Expect(p).ToNot(BeNil())
				Expect(p.Size).To(Equal(uint32(2)))
				Expect(p.Data).ToNot(BeNil())
				Expect(len(p.Data)).To(Equal(2))
				Expect(string(p.Data)).To(Equal("go"))
			})
		})
	})
	Describe("Remove", func() {
		It("can remove", func() {
			// data(data size) | data size(4) | magic(2) |
			bytes := []byte{0x66, 0x65, 0x67, 0x6f, 0x00, 0x00, 0x00, 0x02, 0x21, 0x23}

			fs := afero.NewMemMapFs()
			err := afero.WriteFile(fs, "/source", bytes, 0644)
			Expect(err).To(BeNil())

			source, err := fs.OpenFile("/source", os.O_APPEND|os.O_WRONLY, 0644)
			Expect(err).To(BeNil())
			defer source.Close()

			err = patch.Remove(source)
			Expect(err).To(BeNil())

			p, err := patch.Get(source)
			Expect(err).To(BeNil())
			Expect(p).To(BeNil())

			stat, err := fs.Stat("/source")
			Expect(err).To(BeNil())

			Expect(stat.Size()).To(Equal(int64(2)))
		})
	})
})
