package cache

import (
	"os"
)

type TarFolder struct {
	Name  string
	Index TarIndex
}

type TarIndex map[string]int64

func NewTarFolder(name string) *TarFolder {
	tarFolder := &TarFolder{Name: name, Index: TarIndex{}}
	return tarFolder
}

func (tf *TarFolder) FileName() string {
	return tf.Name + ".tar"
}

func (tf *TarFolder) WriteIndex() error {
	file, err := tf.OpenWrite(0)
	if err != nil {
		return err
	}
	defer file.Close()
	// writer := tar.NewWriter(file)
	// header := tar.Header{
	// 	Name: "______index______",
	// 	Size: int64(0),
	// }
	return nil
}

func (tf *TarFolder) OpenWrite(seek int64) (*os.File, error) {
	file, err := os.OpenFile(tf.FileName(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	if seek > 0 {
		if _, err := file.Seek(seek, 0); err != nil {
			file.Close()
			return nil, err
		}
	}
	return file, nil
}
