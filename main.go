package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"tablemaker/output"
	"tablemaker/tables"
)

func main() {
	var (
		inputFile  = flag.String("input", "", "Input JSON file path (required)")
		outputFile = flag.String("out", "", "Output file path (optional, defaults to stdout)")
		pngOutput  = flag.Bool("png", false, "Generate PNG output instead of text")
	)
	flag.Parse()

	if *inputFile == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s -input <json_file> [-out <output_file>] [-png]\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Read and parse JSON file
	data, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		log.Fatalf("Error reading input file: %v", err)
	}

	var config tables.TableConfig
	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	// Generate table text
	renderer, err := tables.GetRenderer(config.Type)
	if err != nil {
		log.Fatalf("Error getting renderer: %v", err)
	}

	tableText := renderer.Render(config)
	if tableText == "" {
		log.Fatal("Generated table is empty")
	}

	// Handle output
	if *pngOutput {
		if config.PNG == nil {
			log.Fatal("PNG configuration not found in JSON file")
		}

		pngPath := *outputFile
		if pngPath == "" {
			pngPath = "output.png"
		}

		if err := output.GeneratePNG(tableText, *config.PNG, pngPath); err != nil {
			log.Fatalf("Error generating PNG: %v", err)
		}

		fmt.Printf("PNG generated: %s\n", pngPath)
		return
	}

	// Text output
	if *outputFile != "" {
		if err := ioutil.WriteFile(*outputFile, []byte(tableText), 0644); err != nil {
			log.Fatalf("Error writing output file: %v", err)
		}
		fmt.Printf("Output written to: %s\n", *outputFile)
	} else {
		fmt.Print(tableText)
	}
}
