package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime/pprof"

	classifier "github.com/google/licenseclassifier/v2"
	"github.com/google/licenseclassifier/v2/assets"
)

var cf *classifier.Classifier

type LicenseFinding string

func handleAllLicenses() {
	const folder = "./licenses"
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		log.Fatalf("can't ReadDir %q. error: %v", folder, err)
	}
	for _, f := range files {
		lic, err := os.Open(path.Join(folder, f.Name()))
		if err != nil {
			log.Printf("[ERROR] can't open %q: %v", f.Name(), err)
			continue
		}
		Classify(lic)
	}
}

func main() {
	var err error
	cf, err = assets.DefaultClassifier()
	if err != nil {
		panic("assets.DefaultClassifier: " + err.Error())
	}

	// lic, err := os.Open("./licenses/libssl1.1")
	// Classify(lic)
	handleAllLicenses()

	fMem, err := os.Create("mem.profile")
	if err != nil {
		panic("could not create memory profile: " + err.Error())
	}
	defer fMem.Close() // error handling omitted for example
	if err := pprof.WriteHeapProfile(fMem); err != nil {
		panic("could not write memory profile: " + err.Error())
	}
}

// Classify uses a single classifier to detect and classify the license found in a file
func Classify(r io.Reader) {
	// Use 'github.com/google/licenseclassifier' to find licenses
	result, err := cf.MatchFrom(r)
	if err != nil {
		log.Fatalf("unable to match licenses: %v", err)
	}

	var findings []LicenseFinding
	seen := map[string]struct{}{}
	for _, match := range result.Matches {
		if match.Confidence <= 0.9 {
			continue
		}

		if _, ok := seen[match.Name]; !ok {
			findings = append(findings, LicenseFinding(match.Name))
			seen[match.Name] = struct{}{}
		}
	}
}
