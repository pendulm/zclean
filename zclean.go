package main

import (
	"log"
	"os"
)

import "golang.org/x/sys/unix"

const (
	DEALLOC_FLAG = unix.FALLOC_FL_PUNCH_HOLE | unix.FALLOC_FL_KEEP_SIZE
	CHUNK_SIZE   = 2 * 1024 * 1024
)

func trimFile(fpath string) {
	fileInfo, err := os.Stat(fpath)

	if err != nil {
		if os.IsNotExist(err) {
			return
		} else {
			log.Fatal(err)
		}
	}

	if fileInfo.Mode().IsRegular() == false {
		return
	}

	size := fileInfo.Size()

	f, err := os.OpenFile(fpath, os.O_RDWR, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	for seek := int64(0); seek < size; seek += CHUNK_SIZE {

		if err := unix.Fallocate(int(f.Fd()), DEALLOC_FLAG, seek, CHUNK_SIZE); err != nil {
			log.Fatal(err)
		} else {
			log.Printf("seek: %d - %d", seek, seek+CHUNK_SIZE)
		}
	}

}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: zclean <file>")
	}
	trimFile(os.Args[1])
}
