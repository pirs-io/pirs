package commons

import (
	"bytes"
	"github.com/rs/zerolog/log"
	"io"
	"os"
)

func StreamFileFromPipe(r *io.PipeReader, chunkSize int64, chunkCallback func(chunk []byte) error) error {
	buff := make([]byte, chunkSize)
	eof := false
	for {
		n, err := r.Read(buff)
		if err != nil {
			log.Err(err)
			break
		}
		var toSend []byte
		if err == io.EOF || n == 0 {
			// handle remaining bytes
			toSend = bytes.Trim(buff[:n], "\x00")
			eof = true
		} else {
			toSend = bytes.Trim(buff, "\x00")
		}
		if eof && n == 0 {
			break
		}
		err = chunkCallback(toSend)
		if err != nil {
			log.Err(err)
			return err
		}
	}
	return nil
}

func StreamFileToPipe(fd *os.File, chunkSize int64, pipeWriter *io.PipeWriter) (err error) {
	go func() {
		defer func(fd *os.File) {
			err = fd.Close()
			if err != nil {
				log.Err(err)
			}
		}(fd)
		defer func(pipeWriter *io.PipeWriter) {
			err = pipeWriter.Close()
			if err != nil {
				log.Err(err)
			}
		}(pipeWriter)

		buff := make([]byte, chunkSize)
		for {
			var n int
			n, err = fd.Read(buff)
			if err != nil {
				log.Err(err)
			}
			remaining := buff[:n]
			if err == io.EOF {
				_, err = pipeWriter.Write(remaining)
				if err != nil {
					log.Err(err)
				}
				break
			}
			_, err = pipeWriter.Write(buff)
			if err != nil {
				log.Err(err)
			}
		}
	}()
	return err
}
