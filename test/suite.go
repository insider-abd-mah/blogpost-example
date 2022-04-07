package test

import (
	"io/ioutil"
)

const ResourceBasePath = "resources/"

func ReadStubFile(fileName string) []byte {
	req, err := ioutil.ReadFile(getTestsPath() + ResourceBasePath + fileName)

	if err != nil {
		panic(err)
	}

	return req
}

func getTestsPath() string {
	return getBasePath() + "test/"
}

func getBasePath() string {
	return "../../"
}
