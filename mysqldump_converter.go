package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func createFile(filename string) *os.File {
	outputFile, err := os.Create(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failure during file creation: %s\n", err)
		os.Exit(5)
	}
	return outputFile
}

func main() {
	outputPtr := flag.String("output", "/tmp/mysqldump_converter", "Directory for converted files")
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	log.Printf("Writing output to %s\n", *outputPtr)
	err := os.MkdirAll(*outputPtr, 0700)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create %s: %s\n", *outputPtr, err)
		os.Exit(4)
	}

	schemaRegex := regexp.MustCompile("^-- Host: .*    Database: (.*)$")
	tableRegex := regexp.MustCompile("^-- Table structure for table `(.*)`$")
	dataRegex := regexp.MustCompile("^-- Dumping data for table `(.*)`$")
	filterRegex := regexp.MustCompile("^(DROP TABLE IF EXISTS|LOCK TABLES |UNLOCK TABLES|/\\*![0-9]{5} |--|$)")

	fh, err := os.Open(flag.Args()[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open %s for reading: %s\n", flag.Args()[0], err)
		os.Exit(2)
	}
	defer fh.Close()

	scanner := bufio.NewScanner(fh)
	buf := make([]byte, 0, 10*1024*1024)
	scanner.Buffer(buf, 5*1024*1024)
	var schema, filename string
	var outputFile *os.File
	for scanner.Scan() {
		if schemaRegex.MatchString(scanner.Text()) {
			schema = schemaRegex.FindStringSubmatch(scanner.Text())[1]
			log.Printf("Processing schema: %s\n", schema)
		}

		if schema != "" {
			if tableRegex.MatchString(scanner.Text()) {
				table := tableRegex.FindStringSubmatch(scanner.Text())[1]
				log.Printf("Processing table schema for %s.%s\n", schema, table)
				filename = filepath.Join(*outputPtr, fmt.Sprintf("%s.%s-schema.sql", schema, table))
				outputFile = createFile(filename)
			}

			if dataRegex.MatchString(scanner.Text()) {
				table := dataRegex.FindStringSubmatch(scanner.Text())[1]
				log.Printf("Processing table data for %s.%s\n", schema, table)
				filename = filepath.Join(*outputPtr, fmt.Sprintf("%s.%s.sql", schema, table))
				outputFile = createFile(filename)
			}

			if outputFile != nil {
				if !filterRegex.MatchString(scanner.Text()) {
					outputFile.WriteString(scanner.Text())
					outputFile.WriteString("\n")
				}
			}

		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Failure during file read: %s\n", err)
		os.Exit(3)
	}
}
