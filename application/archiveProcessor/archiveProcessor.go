package archiveProcessor

import (
	"bufio"
	"bytes"
	"os"
)

func Process(file *os.File)  ([]string, error){

	var (
		content []string
		part []byte
		prefix bool
		err error
	)

	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 0))
	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			break
		}
		buffer.Write(part)
		if !prefix {
			content = append(content, buffer.String())
			buffer.Reset()
		}
	}

	return content, nil
}