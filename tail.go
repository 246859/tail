package tail

import (
	"bytes"
	"io"
	"os"
	"runtime"
	"slices"
	"unsafe"
)

// Tail returns the last n lines of the file
func Tail(file *os.File, n int) ([]byte, error) {
	tail, _, err := TailAt(file, n, -1)
	return tail, err
}

// TailLines returns bytes slice that last n lines of the file
func TailLines(file *os.File, n int) ([][]byte, error) {
	lines, _, err := TailAtLines(file, n, -1)
	return lines, err
}

// TailAtStringLines returns the string slice that last n lines of the file, and the current offset of the file
func TailAtStringLines(file *os.File, n int, offset int64) ([]string, int64, error) {
	lines, offset, err := TailAtLines(file, n, offset)
	if err != nil {
		return nil, 0, err
	}
	buf := make([]string, len(lines))
	for i, line := range lines {
		buf[i] = unsafe.String(unsafe.SliceData(line), len(line))
	}
	return buf, offset, nil
}

// TailAtLines returns the byte slice that last n lines of the file, and the current offset of the file
func TailAtLines(file *os.File, n int, offset int64) ([][]byte, int64, error) {
	content, offset, err := TailAt(file, n, offset)
	if err != nil {
		return nil, 0, err
	}
	lines := bytes.Split(content, []byte{'\n'})
	return lines, offset, nil
}

// TailAtString returns the string representation that last n lines of the file, and the current offset of the file
func TailAtString(file *os.File, n int, offset int64) (string, int64, error) {
	content, offset, err := TailAt(file, n, offset)
	if err != nil {
		return "", 0, err
	}
	return unsafe.String(unsafe.SliceData(content), len(content)), offset, nil
}

// TailAt returns the last n lines of the file, and the current offset of the file
func TailAt(file *os.File, n int, offset int64) ([]byte, int64, error) {
	stat, err := file.Stat()
	if err != nil {
		return nil, 0, err
	}
	size := stat.Size()

	if offset == 0 {
		offset = -1
	}
	content := make([]byte, 0, n)
	char := make([]byte, 1)

	for n > 0 {
		// read byte one by one
		for offset >= -size {
			_, err := file.Seek(offset, io.SeekEnd)
			if err != nil {
				return nil, 0, err
			}

			_, err = file.Read(char)
			if err != nil {
				return nil, 0, err
			}
			if char[0] == '\n' {
				break
			}
			content = append(content, char[0])
			offset--
		}

		offset--
		// 如果是windows系统，则需要额外移动偏移量
		if runtime.GOOS == "windows" { // CRLF
			offset--
		}

		// 最后一行不加换行符
		if offset >= -size {
			content = append(content, '\n')
		}
		n--

		// 防止seek指针移动到起始位置之外
		if offset < -size {
			offset = 0
			break
		}
	}

	// 最后再整体反转
	slices.Reverse(content)
	return content, offset, err
}
