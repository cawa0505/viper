package main

import (
	"errors"
	"github.com/artonge/go-csv-tag"
	"io/ioutil"
	"os"
	"strings"
)

type vm struct {
	Name string `csv:"name"`
	URL  string `csv:"url"`
}

func listRemoteImagesNames() ([]string, error) {

	vms, err := listRemoteImages()
	if err != nil {
		return nil, err
	}

	images := make([]string, 0)

	for i := range vms {
		images = append(images, vms[i].Name)
	}

	return images, nil
}

func listRemoteImages() ([]vm, error) {

	const filePattern = "remote.viper.csv_"
	var tmp string
	// read images csv from tmp cache
	dir, err := ioutil.ReadDir(os.TempDir())
	for i := range dir {
		if strings.HasPrefix(dir[i].Name(), filePattern) {
			tmp = os.TempDir() + dir[i].Name()
			break
		}
	}

	// if not found, then download the file
	if len(tmp) == 0 {
		tmpFile, err := ioutil.TempFile("", filePattern)
		if err != nil {
			return nil, err
		}
		tmp = tmpFile.Name()
		_, err = execute("wget", "-O", tmp, "https://raw.githubusercontent.com/mhewedy/viper/master/samples/images.csv")
		if err != nil {
			return nil, err
		}
	}

	// parse the file as csv
	var vms []vm
	err = csvtag.Load(csvtag.Config{
		Path: tmp,
		Dest: &vms,
	})

	if err != nil {
		return nil, err
	}

	err = validate(vms)
	if err != nil {
		return nil, err
	}

	return vms, nil
}

// Check name follows <distro>/<version> name
// Check Unique name
func validate(vms []vm) error {
	names := make(map[string]bool)

	for i := range vms {
		// check name
		if len(strings.Split(vms[i].Name, "/")) != 2 {
			return errors.New("name doesn't follow pattern <distro>/<version>")
		}
		// check duplicate
		_, found := names[vms[i].Name]
		if found {
			return errors.New("remote list cannot contains duplicates")
		}

		names[vms[i].Name] = true
	}
	return nil
}
