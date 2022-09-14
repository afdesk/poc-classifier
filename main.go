package main

import (
	"os"
	"runtime/pprof"

	//classifier "github.com/google/licenseclassifier/v2"
	"github.com/google/licenseclassifier/v2/assets"
)

func main()  {
	assets.DefaultClassifier()

	fMem, err := os.Create("mem.profile")
	if err != nil {
		panic("could not create memory profile: " + err.Error())
	}
	defer fMem.Close() // error handling omitted for example
	if err := pprof.WriteHeapProfile(fMem); err != nil {
		panic("could not write memory profile: " + err.Error())
	}
}