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
	Extensions   []string
	Class        string
}

func (fs FileSignature) SigName() string {
	specials := map[string]string{
		" -|!:@,": "_",
		".()":     "",
		"+":       "P",
	}

	name := strings.ToUpper(fs.Description)
	exts := strings.ToUpper(strings.Join(fs.Extensions, "_"))
	for chars, rep := range specials {
		for _, c := range chars {
			name = strings.Replace(name, string(c), rep, -1)
			exts = strings.Replace(exts, string(c), rep, -1)
		}
	}

	sigName := fmt.Sprintf("SIG_%s", name)
	if fs.Extensions != nil {
		sigName = fmt.Sprintf("%s_%s", sigName, exts)
	}

	return sigName
}

func (fs FileSignature) ExtensionsGoString() string {
	return fmt.Sprintf("%#v", fs.Extensions)
}

func (fs FileSignature) HeaderGoString() string {
	return fmt.Sprintf("%#v", fs.Header)
}

func (fs FileSignature) TrailerGoString() string {
	return fmt.Sprintf("%#v", fs.Trailer)
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
func NewFileSignature(desc, header, exts, class, offset, trailer string) (*FileSignature, error) {
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

	var extList []string
	if exts != "" && exts != "(none)" {
		for _, e := range strings.Split(exts, "|") {
			extList = append(extList, strings.ToUpper(e))
		}
	}

	return &FileSignature{
		Description:  desc,
		Header:       headerBytes,
		Extensions:   extList,
		Class:        class,
		HeaderOffset: offsetUint64,
		Trailer:      trailerBytes,
	}, nil
}
