/******************************************************************************
* This utility uses the file signature tagble to find the file type.
*
*	filetype <file_path>
*
*	Copyright (c) 2020 Sreejith Kesavan <sreejithemk@gmail.com>
*
******************************************************************************/

package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/semk/filesigtable-go/sigtable"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: filetype <file_path>")
		os.Exit(1)
	}
	file := os.Args[1]

	f, err := os.OpenFile(file, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	defer f.Close()

	var headerBuf []byte
	var validSigs []*sigtable.FileSignature
	for _, s := range sigtable.FileSignatures {
		headerBuf = make([]byte, len(s.Header))
		_, err := f.ReadAt(headerBuf, s.HeaderOffset)
		if err != nil {
			fmt.Println(err)
			os.Exit(3)
		}
		if bytes.Compare(headerBuf, s.Header) == 0 {
			validSigs = append(validSigs, s)
		}
	}

	if validSigs != nil {
		fmt.Println("File signatures matched with the following types:\n")
		for _, s := range validSigs {
			fmt.Printf("Description:\t%s\n", s.Description)
			fmt.Printf("Class:\t\t%s\n", s.Class)
			fmt.Printf("Extension:\t%s\n\n", s.Extension)
		}
	} else {
		fmt.Println("No matching signatures found.")
		os.Exit(4)
	}

}
