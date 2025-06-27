package infracost

import _ "embed"

//go:embed infracost-usage-example.yml
var referenceUsageFileContents []byte

//go:embed images.json
var imageFileContents []byte

func GetReferenceUsageFileContents() *[]byte {
	return &referenceUsageFileContents
}

func GetImageFileContent() *[]byte {
	return &imageFileContents
}
