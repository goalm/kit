package read

import (
	"bufio"
	"fmt"
	"github.com/goalm/kit/sys"
	"github.com/google/go-cmp/cmp"
	"log"
	"os"
	"strconv"
	"strings"
)

func CsvToEnum(filePath string) *sys.Enum {
	m := sys.NewEnum()
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	rowNo := 0
	for scanner.Scan() {
		rowNo++
		line := scanner.Text()
		if line == "\xA0" || line == "" {
			break
		}
		ss := strings.Split(line, ",")
		if rowNo == 1 {
			validEnumHeader(line)
			continue
		}

		idx, err := strconv.Atoi(ss[0])
		if err != nil {
			log.Fatalf("Error converting %v to int", ss[0])
		}
		name := ss[1]
		desc := ss[2]

		m.Add(idx, name, desc)
	}
	return m
}

func validEnumHeader(line string) {
	ss := strings.Split(line, ",")
	if !cmp.Equal(ss, []string{"Index", "Name", "Description"}) {
		fmt.Println("Invalid header for Enumeration: ", ss, "Expected: Index, Name, Description")
		log.Println("Invalid header for Enumeration: ", ss, "Expected: Index, Name, Description")
	}
}
