package patch

import (
	"encoding/binary"
	"fmt"
	"io/fs"

	"github.com/spf13/afero"
)

type Patch struct {
	Magic []byte
	Size  uint32
	Data  []byte
}

var MagicNumber = [...]byte{0x21, 0x23}

type Reader interface {
	fs.File
}

func FromString(data string) *Patch {
	bytes := []byte(data)
	return &Patch{
		Magic: MagicNumber[:],
		Size:  uint32(len(bytes)),
		Data:  bytes,
	}
}

func (p *Patch) Apply(target afero.File) error {
	_, err := target.Write(p.Data)
	if err != nil {
		return err
	}

	sizeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(sizeBytes, p.Size)

	_, err = target.Write(sizeBytes)
	if err != nil {
		return err
	}

	_, err = target.Write(p.Magic)
	return err
}

func Remove(source afero.File) error {
	dataSize, err := readDataSize(source)
	if err != nil {
		return err
	}
	offset, err := getDataOffset(source, dataSize)
	if err != nil {
		return err
	}

	return source.Truncate(offset)
}

func Get(source afero.File) (*Patch, error) {
	mn, err := readMagicNumber(source)
	if err != nil {
		return nil, nil
	}

	if !validateMagicNumber(mn) {
		return nil, nil
	}

	dataSize, err := readDataSize(source)
	if err != nil {
		return nil, err
	}

	data, err := readData(source, dataSize)
	if err != nil {
		return nil, err
	}

	return &Patch{
		Magic: MagicNumber[:],
		Size:  dataSize,
		Data:  data,
	}, nil
}

func readMagicNumber(source afero.File) ([]byte, error) {
	offset, err := getMagicNumberOffset(source)
	if err != nil {
		return nil, err
	}
	bytes := make([]byte, len(MagicNumber))
	i, err := source.ReadAt(bytes, offset)
	if err != nil {
		return nil, err
	}
	if i != len(MagicNumber) {
		return nil, fmt.Errorf("expected %d bytes read, actual %d", len(bytes), i)
	}
	return bytes, nil
}

func getMagicNumberOffset(source afero.File) (int64, error) {
	stat, err := source.Stat()
	if err != nil {
		return 0, err
	}

	offset := stat.Size() - int64(len(MagicNumber))
	if offset < 0 {
		return 0, fmt.Errorf("offset is less than 0")
	}
	return offset, nil
}

func validateMagicNumber(data []byte) bool {
	for i, b := range MagicNumber {
		if b != data[i] {
			return false
		}
	}
	return true
}

func readDataSize(source afero.File) (uint32, error) {
	offset, err := getDataSizeOffset(source)
	if err != nil {
		return 0, err
	}
	bytes := make([]byte, 4)
	i, err := source.ReadAt(bytes, offset)
	if err != nil {
		return 0, err
	}
	if i != len(bytes) {
		return 0, fmt.Errorf("expected %d bytes read, actual %d", len(bytes), i)
	}

	return binary.BigEndian.Uint32(bytes), nil
}

func getDataSizeOffset(source afero.File) (int64, error) {
	offset, err := getMagicNumberOffset(source)
	if err != nil {
		return 0, err
	}
	return offset - 4, nil
}

func readData(source afero.File, dataSize uint32) ([]byte, error) {
	offset, err := getDataOffset(source, dataSize)
	if err != nil {
		return nil, err
	}
	bytes := make([]byte, dataSize)
	i, err := source.ReadAt(bytes, offset)
	if err != nil {
		return nil, err
	}
	if i != len(bytes) {
		return nil, fmt.Errorf("expected %d bytes read, actual %d", len(bytes), i)
	}
	return bytes, nil
}

func getDataOffset(source afero.File, dataSize uint32) (int64, error) {
	offset, err := getDataSizeOffset(source)
	if err != nil {
		return 0, err
	}
	return offset - int64(dataSize), nil
}
