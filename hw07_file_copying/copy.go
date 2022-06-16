package main

import (
	"errors"
	"io"
	"os"

	pb "github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrOpenFileFrom          = errors.New("can`t open file FROM")
	ErrOpenFileTo            = errors.New("can`t open file TO")
	ErrWrite                 = errors.New("can`t write to file")
	ErrSeek                  = errors.New("can`t change position cursor")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fileFrom, errOpen := os.Open(fromPath)

	if errOpen != nil {
		return ErrOpenFileFrom
	}

	defer fileFrom.Close()

	fileInfo, errStat := fileFrom.Stat()

	if errStat != nil {
		return ErrOpenFileFrom
	}

	if fileInfo.Size() <= 0 {
		return ErrUnsupportedFile
	}

	if fileInfo.Size() < offset {
		return ErrOffsetExceedsFileSize
	}

	if limit > fileInfo.Size() {
		limit = fileInfo.Size()
	}

	_, errSeek := fileFrom.Seek(offset, 0)

	if errSeek != nil {
		return ErrSeek
	}

	fileTo, errCreate := os.Create(toPath)

	if errCreate != nil {
		return ErrOpenFileTo
	}

	barSize := fileInfo.Size() - offset
	if limit > 0 {
		barSize = limit
	}

	bar := pb.Full.Start64(barSize)
	varFileFrom := bar.NewProxyReader(fileFrom)

	var errWrite error

	if limit > 0 {
		_, errWrite = io.CopyN(fileTo, varFileFrom, limit)
	} else {
		_, errWrite = io.Copy(fileTo, varFileFrom)
	}

	if errWrite != nil {
		return ErrWrite
	}

	bar.Finish()

	return nil
}
