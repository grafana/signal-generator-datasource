// Transform Influx format file.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

// go run cmd/transform/transform.go -sort-fields -drop-field hv_dc_voltage_V -drop-field valid dev/data.log

type arrayFlags []string

func (i *arrayFlags) String() string {
	return strings.Join(*i, ", ")
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func main() {
	var dropFields arrayFlags
	var keepFields arrayFlags
	var sortFields bool
	var nanToZero bool
	var dropNan bool

	printUsage := func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "transform [flags] file\n")
		flag.PrintDefaults()
	}

	flag.Usage = func() {
		printUsage()
	}

	flag.Var(&dropFields, "drop-field", "Drop field")
	flag.Var(&keepFields, "keep-field", "Keep field")
	flag.BoolVar(&sortFields, "sort-fields", false, "Sort fields")
	flag.BoolVar(&nanToZero, "nan-to-zero", false, "Replace NaN with 0")
	flag.BoolVar(&dropNan, "drop-nan", false, "Drop NaN fields")
	flag.Parse()

	if flag.NArg() == 0 {
		printUsage()
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "\ninput file required\n")
		os.Exit(1)
	}

	path := flag.Args()[0]

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		fieldSet := strings.Split(fields[1], ",")

		if len(dropFields) > 0 {
			var newFieldSet []string
			for _, field := range fieldSet {
				fieldName := strings.Split(field, "=")[0]
				if contains(dropFields, fieldName) {
					continue
				}
				newFieldSet = append(newFieldSet, field)
			}
			fieldSet = newFieldSet
		}

		if len(keepFields) > 0 {
			var newFieldSet []string
			for _, field := range fieldSet {
				fieldName := strings.Split(field, "=")[0]
				if contains(keepFields, fieldName) {
					newFieldSet = append(newFieldSet, field)
				}
			}
			fieldSet = newFieldSet
		}

		if nanToZero || dropNan {
			var newFieldSet []string
			for _, field := range fieldSet {
				fieldName := strings.Split(field, "=")[0]
				fieldValue := strings.Split(field, "=")[1]
				if fieldValue == `"NaN"` {
					if nanToZero {
						newFieldSet = append(newFieldSet, fmt.Sprintf("%s=%d", fieldName, 0))
					} else {
						continue
					}
				} else {
					newFieldSet = append(newFieldSet, field)
				}
			}
			fieldSet = newFieldSet
		}

		if sortFields {
			sort.Slice(fieldSet, func(i, j int) bool {
				return fieldSet[i] < fieldSet[j]
			})
		}

		fields[1] = strings.Join(fieldSet, ",")

		if _, err = os.Stdout.WriteString(strings.Join(fields, " ") + "\n"); err != nil {
			log.Fatal(err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
