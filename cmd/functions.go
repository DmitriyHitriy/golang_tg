package functions

import (
	"bufio"
	"os"
)

func Get_rows_in_file(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	buf := bufio.NewScanner(file)
	rows := make([]string, 0)

	for buf.Scan() {
		rows = append(rows, buf.Text())
	}

	return rows, nil
}
