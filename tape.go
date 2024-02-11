package poker

import "os"

type tape struct {
	file *os.File
}

// custom write for the fs db file, makes sure the file is clear to write
func (t *tape) Write(p []byte) (n int, err error) {
	// clear the file before writing the league to it
	t.file.Truncate(0)
	t.file.Seek(0, 0)
	return t.file.Write(p)
}
