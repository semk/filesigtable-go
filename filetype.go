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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/semk/filesigtable-go/sigtable"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: filetype <file_path>")
		os.Exit(1)
	}
	file := os.Args[1]

	ext := strings.ToUpper(filepath.Ext(file))
	f, err := os.OpenFile(file, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	defer f.Close()

	var validSigs []sigtable.FileSignature

	for _, s := range sigtable.GetSignaturesByExtension(ext) {
		valid, err := sigtable.ValidateSignature(s, f)
		if err != nil {
			fmt.Println(err)
			os.Exit(3)
		}
		if valid {
			validSigs = append(validSigs, s)
		}
	}

	if validSigs != nil {
		fmt.Println("File signatures matched with the following types:\n")
		for _, s := range validSigs {
			fmt.Printf("Description:\t%s\n", s.Description)
			fmt.Printf("Class:\t\t%s\n", s.Class)
			fmt.Printf("Extensions:\t%s\n\n", strings.Join(s.Extensions, ", "))
		}
	} else {
		fmt.Println("No matching signatures found by extension. Searching greedily!")
		for _, s := range sigtable.GetAllSignatures() {
			valid, err := sigtable.ValidateSignature(s, f)
			if err != nil {
				fmt.Println(err)
				os.Exit(3)
			}
			if valid {
				validSigs = append(validSigs, s)
			}
		}
		if validSigs != nil {
			fmt.Println("File signatures matched with the following types:\n")
			for _, s := range validSigs {
				fmt.Printf("Description:\t%s\n", s.Description)
				fmt.Printf("Class:\t\t%s\n", s.Class)
				fmt.Printf("Extensions:\t%s\n\n", strings.Join(s.Extensions, ", "))
			}
		} else {
			fmt.Println("No matching signatures found.")
			os.Exit(4)
		}
	}

}
