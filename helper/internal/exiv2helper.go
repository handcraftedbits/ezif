package internal // import "golang.handcraftedbits.com/ezif/helper/internal"

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"os/exec"
)

//
// Public types
//

type ExternalExiv2 interface {
	Add(key string, values []interface{}) ExternalExiv2
	Set(key string, values []interface{}) ExternalExiv2
}

//
// Private types
//

// ExternalExiv2 implementation
type externalExiv2Impl struct {
	args         []string
	image        []byte
	tempFilename string
}

func (exiv2 *externalExiv2Impl) Add(name string, values []interface{}) ExternalExiv2 {
	exiv2.args = append(exiv2.args, "-M", fmt.Sprintf("add %s %s", name, convertValuesToExiv2Format(values)))

	return exiv2
}

func (exiv2 *externalExiv2Impl) Set(name string, values []interface{}) ExternalExiv2 {
	exiv2.args = append(exiv2.args, "-M", fmt.Sprintf("set %s %s", name, convertValuesToExiv2Format(values)))

	return exiv2
}

func (exiv2 *externalExiv2Impl) execute() (error, string, string) {
	var args []string
	var command *exec.Cmd
	var err error
	var imageFilename string
	var stdErr bytes.Buffer
	var stdOut bytes.Buffer

	imageFilename, err = saveImage(exiv2.image)

	if err != nil {
		return err, "", ""
	}

	exiv2.tempFilename = imageFilename

	args = append(exiv2.args, imageFilename)
	command = exec.Command("exiv2", args...)

	command.Stderr = &stdErr
	command.Stdout = &stdOut

	if err := command.Run(); err != nil {
		return err, stdOut.String(), stdErr.String()
	}

	return nil, stdOut.String(), stdErr.String()
}

//
// Private functions
//

func convertValuesToExiv2Format(values []interface{}) string {
	var buffer bytes.Buffer

	for i, value := range values {
		switch v := value.(type) {
		case *big.Rat:
			buffer.WriteString(fmt.Sprintf("%v/%v", v.Num(), v.Denom()))
		case xmpLangAltEntry:
			buffer.WriteString(fmt.Sprintf("lang=\"%s\" %s", v.language, v.value))
		default:
			buffer.WriteString(fmt.Sprintf("%v", v))
		}

		if i < len(values)-1 {
			buffer.WriteRune(' ')
		}
	}

	return buffer.String()
}

func newExternalExiv2(image []byte) *externalExiv2Impl {
	return &externalExiv2Impl{
		image: image,
	}
}

func saveImage(image []byte) (string, error) {
	var err error
	var tempFile *os.File

	tempFile, err = ioutil.TempFile("", "ezif-test")

	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(tempFile.Name(), image, os.ModePerm)

	if err != nil {
		return "", err
	}

	_ = tempFile.Close()

	return tempFile.Name(), nil
}
