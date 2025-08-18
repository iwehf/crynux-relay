package blockchain

import (
	"crypto/sha256"
	"encoding/binary"
	"image/png"
	"io"

	"github.com/corona10/goimagehash"
)


func GetPHashForImage(reader io.Reader) ([]byte, error) {
	image, err := png.Decode(reader)
	if err != nil {
		return nil, err
	}
	pHash, err := goimagehash.PerceptionHash(image)
	if err != nil {
		return nil, err
	}

	bs := make([]byte, pHash.Bits()/8)
	binary.BigEndian.PutUint64(bs, pHash.GetHash())
	return bs, nil
}

func GetHashForGPTResponse(reader io.Reader) ([]byte, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	h := sha256.Sum256(content)
	return h[:], nil
}
