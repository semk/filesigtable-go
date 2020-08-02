/******************************************************************************
* Structures and utility functions to work with file signatures from Kessler's
* data file.
*
*	Copyright (c) 2020 Sreejith Kesavan <sreejithemk@gmail.com>
*
******************************************************************************/

package sigtable

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

// FileSignature represents information about a file's signature.
// Each comma separated value from Gary Kessler's FileSig (RAW)
type FileSignature struct {
	Description  string
	Header       []byte
	HeaderOffset int64
	HeaderLength uint64
	Trailer      []byte
	Extension    string
	Class        string
}

func (fs FileSignature) HeaderGoString() string {
	if fs.Header == nil {
		return "nil"
	}
	var hexBytes []string
	for _, b := range fs.Header {
		hexByte := fmt.Sprintf("0x%x", b)
		hexBytes = append(hexBytes, hexByte)
	}
	return "[]byte{" + strings.Join(hexBytes, ", ") + "}"
}

func (fs FileSignature) TrailerGoString() string {
	if fs.Trailer == nil {
		return "nil"
	}
	var hexBytes []string
	for _, b := range fs.Trailer {
		hexByte := fmt.Sprintf("0x%x", b)
		hexBytes = append(hexBytes, hexByte)
	}
	return "[]byte{" + strings.Join(hexBytes, ", ") + "}"
}

func convertHexStringToBytes(hexStr string) ([]byte, error) {
	var result []byte
	for _, s := range strings.Split(hexStr, " ") {
		if strings.HasPrefix(s, "?") {
			// FIXME: Hanle wildcard bytes
			continue
		}
		b, err := hex.DecodeString(s)
		if err != nil {
			return nil, err
		}
		result = append(result, b...)
	}
	return result, nil
}

// NewFileSignature creates a new FileSignature object from string representation
func NewFileSignature(desc, header, ext, class, offset, trailer string) (*FileSignature, error) {
	var err error
	var headerBytes []byte
	if header != "(null)" {
		headerBytes, err = convertHexStringToBytes(header)
		if err != nil {
			return nil, err
		}
	}

	var trailerBytes []byte
	if trailer != "(null)" {
		trailerBytes, err = convertHexStringToBytes(trailer)
		if err != nil {
			return nil, err
		}
	}

	offsetUint64, err := strconv.ParseInt(offset, 10, 64)
	if err != nil {
		return nil, err
	}

	return &FileSignature{
		Description:  desc,
		Header:       headerBytes,
		Extension:    ext,
		Class:        class,
		HeaderOffset: offsetUint64,
		Trailer:      trailerBytes,
	}, nil
}
