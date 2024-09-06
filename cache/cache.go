package cache

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/VictoriaMetrics/fastcache"
)

type Table struct {
	Caches  *fastcache.Cache // cache for each table, with rowKey as keys
	NumRows int              // number of rows
	ColKeys []string         // sub key for each record
}

func LoadGenericTable(filePath string, maxBytes int) *Table {
	start := time.Now()
	cache := fastcache.New(maxBytes)
	colKeys := make([]string, 0)

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	noIdx := 0
	scanner := bufio.NewScanner(file)
	numRows := 0
	for scanner.Scan() {
		line := scanner.Text()
		// end of file
		if line == "\xA0" || line == "" {
			break
		}
		// dump descriptions
		if line[0] != '!' && line[0] != '*' {
			continue
		}

		line = strings.ReplaceAll(line, "\"", "")
		// process header
		if line[0] == '!' {
			noIdx, err = strconv.Atoi(line[1:2])
			if err != nil {
				log.Fatal(err)
			}
			if noIdx < 1 {
				log.Fatalf("Table %v has no keys: %v", filePath, line)
			}
			str := strings.Split(line, ",")
			colKeys = str[noIdx:]
			cache.Set([]byte("subKeys"), []byte(strings.Join(colKeys, ",")))
			// process data
		} else if line[0] == '*' {
			numRows++
			str := strings.Split(line, ",")
			rowKeys := str[1:noIdx]

			key := rowKeys[0]
			for _, v := range rowKeys[1:] {
				key = key + ":" + v
			}

			colVals := []byte(strings.Join(str[noIdx:], ","))
			cache.Set([]byte(key), colVals)
		}
	}
	fmt.Printf("loading data %s, total rows%d, used %v\n", filePath, numRows, time.Since(start))

	return &Table{cache, numRows, colKeys}
}
