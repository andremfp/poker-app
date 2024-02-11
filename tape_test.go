package poker_test

import (
	"io"
	"os"
	"testing"
)

type SpyTape struct {
	file *os.File
}

// custom write for the fs db file, makes sure the file is clear to write
func (t *SpyTape) Write(p []byte) (n int, err error) {
	// clear the file before writing the league to it
	t.file.Truncate(0)
	t.file.Seek(0, 0)
	return t.file.Write(p)
}

func TestTapeWrite(t *testing.T) {
	file, clean := createTempFile(t, "12345")
	defer clean()

	tape := &SpyTape{file}

	tape.Write([]byte("abc"))

	file.Seek(0, 0)
	newFileContents, _ := io.ReadAll(file)

	got := string(newFileContents)
	want := "abc"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
