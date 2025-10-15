package fs

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/mws-cloud-platform/util-toolset/pkg/utils/consterr"
)

type stdoutPrint struct {
	FS
}

// WithStdoutPrint is an option for NewFS that makes [WriteOnlyFS] WriteFile
// print its content to stdout instead of writing to a file.
func WithStdoutPrint() Option {
	return func(fs FS) FS {
		return &stdoutPrint{
			FS: fs,
		}
	}
}

const errNotSupported = consterr.Error("stdoutPrint: OpenFile not supported")

func (*stdoutPrint) OpenFile(string, int, os.FileMode) (WritableFile, error) {
	return nil, errNotSupported
}

func (*stdoutPrint) WriteFile(_ string, data []byte, _ fs.FileMode) error {
	_, err := fmt.Fprint(os.Stdout, string(data))
	return err
}
