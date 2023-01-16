package cache_test

import (
	"archive/tar"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
)

const contentsEntryA = "adasdflk;jkj;lkj;alsdfasdfasdjlfkj;dsfsadjfk"
const contentsEntryB = "assssasdfasdfsdffd"

func TestTarFileSystem(t *testing.T) {
	fileName := uuid.NewString()
	file, err := os.Create(fileName)
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		file.Close()
		os.Remove(fileName)
	}()
	writer := tar.NewWriter(file)
	header := tar.Header{
		Name: "one",
		Size: int64(len([]byte(contentsEntryA))),
	}
	if err := writer.WriteHeader(&header); err != nil {
		t.Error(err)
		return
	}
	writer.Write([]byte(contentsEntryA))
	if err := writer.Flush(); err != nil {
		t.Error(err)
		return
	}
	posAfterFirstEntry, err := file.Seek(0, 1)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("posAfterFirstEntry", posAfterFirstEntry)
	header = tar.Header{
		Name: "two",
		Size: int64(len([]byte(contentsEntryB))),
	}
	if err := writer.WriteHeader(&header); err != nil {
		t.Error(err)
		return
	}
	writer.Write([]byte(contentsEntryB))
	posAfterSecondEntry, err := file.Seek(0, 1)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("posAfterSecondEntry", posAfterSecondEntry)
	file.Close()

	infile, err := os.Open(fileName)
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		infile.Close()
	}()
	fmt.Println("Seek", posAfterFirstEntry)
	if _, err = infile.Seek(posAfterFirstEntry, 0); err != nil {
		t.Error(err)
		return
	}
	reader := tar.NewReader(infile)
	inheader, err := reader.Next()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(inheader.Name)

}
