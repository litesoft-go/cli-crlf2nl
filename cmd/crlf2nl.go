package main

import (
    "lib-builtin/lib/fatal"
    "lib-builtin/lib/cli"

    "strings"
    "lib-builtin/lib/validate"
    "github.com/pkg/errors"
    "os"
    "path/filepath"
    "io/ioutil"
    "bytes"
)

var Extension string

func visitor(path string, f os.FileInfo, err error) error {
    if (err == nil) && !f.IsDir() && strings.HasSuffix(path, Extension) {
	err = process(path, f.Mode())
    }
    return err
}

func main() {
    zCLI := cli.New().ShowVersion("1.0").
	    AddRequiredString(&Extension, "Extension", "Extention to process",
	validate.NewStringValidationBuilder().AddMutator(func(pValue string) (rValidated string, err error) {
	    rValidated = strings.TrimSpace(pValue)
	    if len(rValidated) == 0 {
		err = errors.New("empty")
	    }
	    return
	}).AsMutator())

    fatal.IfErr(zCLI.Parse())
    println("Extension:", Extension)
    if !strings.HasPrefix(Extension, ".") {
	Extension = "." + Extension
    }

    zWorkingDir, err := os.Getwd()
    fatal.IfErr(err)
    fatal.IfErr(filepath.Walk(zWorkingDir, visitor))
}

func process(pFilePath string, pCurrentMode os.FileMode) (err error) {
    print(" ", pFilePath)
    zReadBytes, err := ioutil.ReadFile(pFilePath)
    if err != nil {
	println(" |", err)
	return
    }
    zNewBytes := bytes.Replace(zReadBytes, []byte{13, 10}, []byte{10}, -1)
    if len(zReadBytes) == len(zNewBytes) {
	println(" : No Change")
	return
    }

    zNewFileName := pFilePath + "-New" + Extension
    err = ioutil.WriteFile(zNewFileName, zNewBytes, pCurrentMode)
    if err == nil {
	err = os.Remove(pFilePath)
    }
    if err == nil {
	err = os.Rename(zNewFileName, pFilePath)
    }
    if err != nil {
	println(" |", err)
    } else {
	println(" : Updated")
    }
    return
}
