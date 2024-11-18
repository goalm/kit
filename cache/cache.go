package cache

import (
	"bufio"
	"fmt"
	"github.com/samber/lo"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/VictoriaMetrics/fastcache"
)

type Table struct {
	Path    string
	Caches  *fastcache.Cache // cache for each table, with rowKey as keys
	NumIdx  int              // number of keys
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
			//cache.Set([]byte("ColKeys"), []byte(strings.Join(colKeys, ",")))
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
	fmt.Printf("loading data %s, in total %d rows, used  %v\n", filePath, numRows, time.Since(start))

	return &Table{filePath, cache, noIdx, numRows, colKeys}
}

func (t *Table) TEXT(err string, idx ...string) string {
	return ReadCacheTableString(t, err, idx...)
}
func (t *Table) INT(err string, idx ...string) int {
	return ReadCacheTableInt(t, err, idx...)
}
func (t *Table) NUM(err string, idx ...string) float64 {
	return ReadCacheTableFloat(t, err, idx...)
}

func ReadCacheTableString(table *Table, error string, idx ...string) string {
	// get main key
	key := idx[0]
	for _, v := range idx[1 : len(idx)-1] {
		key = key + ":" + v
	}
	// validate key number
	if table.NumIdx != len(idx) {
		log.Fatalf("Key number not match for %s, expected:%d \n", table, table.NumIdx)
		return ""
	}

	rowVal := table.Caches.Get(nil, []byte(key))

	// get column key
	n := lo.IndexOf(table.ColKeys, idx[len(idx)-1])
	// record does not match keys
	if n == -1 {
		if error == "Y" || error == "y" {
			log.Fatalf("Key %s not found in %s\n", idx[len(idx)-1], table)
		}
		if error == "N" || error == "n" {
			return ""
		}
	}

	rowStr := strings.Split(string(rowVal), ",")
	if len(rowStr) < n+1 {
		return ""
	}
	return rowStr[n]
}

func ReadCacheTableInt(table *Table, error string, idx ...string) int {
	str := ReadCacheTableString(table, error, idx...)
	if str == "" {
		return 0
	}
	val, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}
	return val
}

func ReadCacheTableFloat(table *Table, error string, idx ...string) float64 {
	str := ReadCacheTableString(table, error, idx...)
	if str == "" {
		return 0
	}
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		panic(err)
	}
	return val
}
