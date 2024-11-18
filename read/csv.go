package read

import (
	"bufio"
	"encoding/csv"
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/goalm/kit/data"
	"github.com/goalm/kit/sys"
	"github.com/google/go-cmp/cmp"
	"github.com/jszwec/csvutil"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

// CsvToEnum read enumeration from CSV file
func CsvToEnum(filePath string) *sys.Enum {
	m := sys.NewEnum()
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	rowNo := 0

	for {
		peek, err := reader.Peek(1)
		if err != nil {
			if err.Error() == "EOF" {
				log.Println("EOF")
				break
			}
		}
		if len(peek) == 0 {
			break
		}

		line, err := reader.ReadString('\n')
		// remove space, ZWNBSP
		line = strings.TrimSpace(line)
		line = strings.Replace(line, "\ufeff", "", -1)

		ss := strings.Split(line, ",")
		if rowNo == 0 {
			validEnumHeader(line)
		} else {
			idx, err := strconv.Atoi(ss[0])
			if err != nil {
				log.Fatalf("Error converting %v to int", ss[0])
			}
			name := ss[1]
			desc := ss[2]
			m.Add(idx, name, desc)
		}
		rowNo++
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

func CsvToProductDimensions(filePath string) map[string]map[string]string {
	m := make(map[string]map[string]string)
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	csvReader := csv.NewReader(reader)
	dec, err := csvutil.NewDecoder(csvReader)
	if err != nil {
		log.Printf("Error occured reading %s, %v", filePath, err.Error())
		return nil
	}

	header := dec.Header()
	for {
		record := &data.ProductProperties{Dimensions: make(map[string]string)}
		if err := dec.Decode(record); err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Error occured reading %s, %v", filePath, err.Error())
			log.Fatal(err)
		}

		for _, i := range dec.Unused() {
			record.Dimensions[header[i]] = dec.Record()[i]
		}

		if _, ok := m[record.ProdName]; !ok {
			m[record.ProdName] = record.Dimensions
		} else {
			log.Fatalf("Duplicated product name: %v in table %s", record.ProdName, filePath)
		}

	}

	if len(m) == 0 {
		log.Fatalf("No record found in table %s", filePath)
	}

	return m
}

func CsvToVariableList(filePath string) []*data.VariableProperties {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	csvReader := csv.NewReader(reader)
	dec, err := csvutil.NewDecoder(csvReader)
	if err != nil {
		log.Printf("Error occured reading %s, %v", filePath, err.Error())
		return nil
	}

	header := dec.Header()
	var variables []*data.VariableProperties
	for {
		record := &data.VariableProperties{OtherData: make(map[string]string)}
		if err := dec.Decode(record); err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Error occured reading %s, %v", filePath, err.Error())
			log.Fatal(err)
		}

		for _, i := range dec.Unused() {
			record.OtherData[header[i]] = dec.Record()[i]
		}

		variables = append(variables, record)
	}

	return variables
}

func CsvToResultList(filePath string) []*data.ResultSpecs {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	csvReader := csv.NewReader(reader)
	dec, err := csvutil.NewDecoder(csvReader)
	if err != nil {
		log.Printf("Error occured reading %s, %v", filePath, err.Error())
		return nil
	}

	header := dec.Header()
	var results []*data.ResultSpecs
	for {
		record := &data.ResultSpecs{OtherData: make(map[string]string)}
		if err := dec.Decode(record); err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Error occured reading %s, %v", filePath, err.Error())
			log.Fatal(err)
		}

		for _, i := range dec.Unused() {
			record.OtherData[header[i]] = dec.Record()[i]
		}

		results = append(results, record)
	}

	return results
}

func CsvToPathTrack(filePath string) []*data.PathTrack {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	csvReader := csv.NewReader(reader)
	dec, err := csvutil.NewDecoder(csvReader)
	if err != nil {
		log.Printf("Error occured reading %s, %v", filePath, err.Error())
		return nil
	}

	header := dec.Header()
	var pathTracks []*data.PathTrack
	for {
		record := &data.PathTrack{OtherData: make(map[string]string)}
		if err := dec.Decode(record); err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Error occured reading %s, %v", filePath, err.Error())
			log.Fatal(err)
		}

		for _, i := range dec.Unused() {
			record.OtherData[header[i]] = dec.Record()[i]
		}

		pathTracks = append(pathTracks, record)
	}

	return pathTracks
}

// Streaming data from CSV file
// - Model point
func StreamModelPoint[T any](filePath string, row T, dataChn chan *T) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	// Skip the header lines
	l := 0
	for {
		peek, err := reader.Peek(1)
		if err != nil {
			if err.Error() == "EOF" {
				log.Printf("Reading completed for %s, %s lines in total, no MP data found", filePath, strconv.Itoa(l))
				return
			}
		}

		l++
		if string(peek) != "!" && string(peek) != "*" {
			// Skip the line
			_, err := reader.ReadString('\n')
			if err != nil {
				if err.Error() == "EOF" {
					log.Printf("Reading completed for %s, %s lines in total (last line w/o line break), no MP data found", filePath, strconv.Itoa(l))
					return
				}
			}
		} else {
			break
		}
	}
	csvReader := csv.NewReader(reader)
	dec, err := csvutil.NewDecoder(csvReader)
	if err != nil {
		log.Printf("Error occured reading %s, %v", filePath, err.Error())
		return
	}

	j := 0
	for {
		j++
		peek, err := reader.Peek(1)
		if err != nil {
			if err.Error() == "EOF" {
				return
			}
		}

		if string(peek) != "!" && string(peek) != "*" {
			// Skip the line
			return
		}

		record := row
		err = dec.Decode(&record)
		if err != nil {
			log.Printf("Error occured reading %s, %v", filePath, err.Error())
			break
		}
		dataChn <- &record
	}
}

func CsvKeySet(filePath string) mapset.Set[string] {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	set := mapset.NewSet[string]()
	scanner := bufio.NewScanner(file)
	noIdx := 0
	//rowKeys := make([]string, 0)
	rowKeysHeader := make([]string, 0)
	colKeys := make([]string, 0)
	//keySlice := make([]string, 0)

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
		if line[0] == '!' {
			noIdx, err = strconv.Atoi(line[1:2])
			if err != nil {
				log.Fatal(err)
			}
			if noIdx < 1 {
				log.Fatalf("Table %v has no keys: %v", filePath, line)
			}
			str := strings.Split(line, ",")
			rowKeysHeader = str[1:noIdx]
			colKeys = str[noIdx:]

		} else if line[0] == '*' {
			str := strings.Split(line, ",")
			rowKeys := str[1:noIdx]

			for i, _ := range rowKeys {
				rowKeys[i] = rowKeysHeader[i] + "-" + rowKeys[i]
			}

			for _, v := range colKeys {
				keySlice := make([]string, len(rowKeys))
				// use copy to avoid modifying the original slice
				copy(keySlice, rowKeys)
				keySlice = append(keySlice, v)
				sort.Strings(keySlice)
				key := strings.Join(keySlice, ":")
				set.Add(key)
			}
		}
	}

	return set
}
