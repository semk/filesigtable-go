// This utility generates the file signature Go file from
// Gary Kessler's excellent FileSigs bundle
// https://www.garykessler.net/software/index.html#filesigs

//go:generate go run generate.go -f ../sigtable/data/file_sigs_RAW.txt -o ../sigtable/sigtable.go

package main

import (
	"encoding/csv"
	"flag"
	"io"
	"log"
	"os"
	"text/template"

	"github.com/semk/filesigtable-go/sigtable"
)

func main() {
	signatureFile := flag.String("f", "file_sigs_RAW.txt", "Path to Gary Kessler's FileSigs RAW file.")
	outFile := flag.String("o", "sigtable.go", "Go definition to be generated from the RAW file.")
	flag.Parse()

	csvFile, err := os.Open(*signatureFile)
	if err != nil {
		log.Fatalln("Couldn't open the signature file", err)
	}

	type templateData struct {
		Sigs      []*sigtable.FileSignature
		SigsByExt map[string][]*sigtable.FileSignature
	}

	var sigs []*sigtable.FileSignature
	sigsByExt := make(map[string][]*sigtable.FileSignature)
	r := csv.NewReader(csvFile)

	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		desc := record[0]
		header := record[1]
		extension := record[2]
		class := record[3]
		offset := record[4]
		trailer := record[5]

		s, err := sigtable.NewFileSignature(desc, header, extension, class, offset, trailer)
		if err != nil {
			log.Fatalln("Couldn't parse the file signature", err)
		}
		for _, ext := range s.Extensions {
			sigsByExt[ext] = append(sigsByExt[ext], s)
		}
		sigs = append(sigs, s)
	}

	tmplData := &templateData{Sigs: sigs, SigsByExt: sigsByExt}
	sigTmpl, err := template.New("signatures.gotmpl").ParseFiles("signatures.gotmpl")
	if err != nil {
		log.Fatalln("Couldn't parse the template file", err)
	}

	out, err := os.OpenFile(*outFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	err = sigTmpl.Execute(out, tmplData)
	if err != nil {
		log.Fatal(err)
	}
}
