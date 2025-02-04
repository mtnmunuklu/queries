package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Rule represents the structure of the YML file to extract the ID.
type Rule struct {
	ID string `yaml:"id"`
}

// RunCommand executes the given command and returns the JSON output.
func RunCommand(alterixPath, inputFilePath, configPath string) ([]byte, error) {
	cmd := exec.Command(alterixPath, "-sigma", "-filepath", inputFilePath, "-config", configPath, "-json", "-cs")
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	err := cmd.Run()
	if err != nil {
		// Log error directly to the error log file
		logErrorToFile(inputFilePath, errOut.String())
		return nil, fmt.Errorf("command execution failed: %v, stderr: %s", err, errOut.String())
	}

	// Check if the output is empty or invalid
	if len(out.Bytes()) == 0 {
		logErrorToFile(inputFilePath, "empty output from command")
		return nil, fmt.Errorf("empty output from command for file %v", inputFilePath)
	}

	return out.Bytes(), nil
}

// ProcessFile processes files in the directory and transforms them into the desired JSON format.
func ProcessFile(alterixPath, inputPath, outputPath, configPath, errorLogPath string) {
	err := filepath.Walk(inputPath, func(inputFilePath string, info os.FileInfo, err error) error {
		if err != nil {
			logErrorToFile(inputFilePath, fmt.Sprintf("failed to read file: %v", err))
			return nil
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Process only YML files
		if filepath.Ext(inputFilePath) == ".yml" {
			// Read ID from YML file
			var rule Rule
			yamlData, err := os.ReadFile(inputFilePath)
			if err != nil {
				logErrorToFile(inputFilePath, fmt.Sprintf("failed to read YML file: %v", err))
				return nil
			}
			err = yaml.Unmarshal(yamlData, &rule)
			if err != nil {
				logErrorToFile(inputFilePath, fmt.Sprintf("failed to unmarshal YML file: %v", err))
				return nil
			}

			// Execute command and get output
			jsonOutput, err := RunCommand(alterixPath, inputFilePath, configPath)
			if err != nil {
				// Already logged in RunCommand, skip further processing
				return nil
			}

			// Check if the output is valid JSON
			var jsonData map[string]interface{}
			if !isValidJSON(jsonOutput) {
				logErrorToFile(inputFilePath, fmt.Sprintf("invalid JSON output: %s", string(jsonOutput)))
				return nil
			}

			// Unmarshal the valid JSON output
			err = json.Unmarshal(jsonOutput, &jsonData)
			if err != nil {
				logErrorToFile(inputFilePath, fmt.Sprintf("failed to unmarshal JSON output: %v", err))
				return nil
			}

			// Integrate ID and JSON data
			jsonData["ID"] = rule.ID
			jsonData["Version"] = "0.1"

			// Determine output file path with .json extension
			outputFilePath := filepath.Join(outputPath, inputFilePath[len(inputPath):])
			outputFilePath = outputFilePath[:len(outputFilePath)-len(filepath.Ext(outputFilePath))] + ".json"

			// Save the transformed JSON to the output file
			outputDir := filepath.Dir(outputFilePath)
			err = os.MkdirAll(outputDir, os.ModePerm)
			if err != nil {
				logErrorToFile(inputFilePath, fmt.Sprintf("failed to create output directory: %v", err))
				return nil
			}

			outputJSON, err := json.MarshalIndent(jsonData, "", "  ")
			if err != nil {
				logErrorToFile(inputFilePath, fmt.Sprintf("failed to marshal JSON: %v", err))
				return nil
			}

			err = os.WriteFile(outputFilePath, outputJSON, 0644)
			if err != nil {
				logErrorToFile(inputFilePath, fmt.Sprintf("failed to write output file: %v", err))
				return nil
			}

			fmt.Printf("Processed and saved output for: %v\n", inputFilePath)
		}

		return nil
	})

	if err != nil {
		logErrorToFile("", fmt.Sprintf("error walking the input path: %v", err))
	}
}

// isValidJSON checks if the given byte slice is a valid JSON.
func isValidJSON(data []byte) bool {
	var js map[string]interface{}
	return json.Unmarshal(data, &js) == nil
}

// logErrorToFile writes error details to the error log file.
func logErrorToFile(inputFilePath, errorMessage string) {
	errorLogFile, err := os.OpenFile("errors.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Failed to open error log file: %v\n", err)
		return
	}
	defer errorLogFile.Close()

	// Log the error with input file path
	if inputFilePath != "" {
		_, err := errorLogFile.WriteString(fmt.Sprintf("Error for file %v: %s\n", inputFilePath, errorMessage))
		if err != nil {
			fmt.Printf("Failed to write to log file: %v\n", err)
		}
	} else {
		// For general errors (e.g., walking path)
		_, err := errorLogFile.WriteString(fmt.Sprintf("General error: %s\n", errorMessage))
		if err != nil {
			fmt.Printf("Failed to write to log file: %v\n", err)
		}
	}
}

func main() {
	// Define command-line flags
	alterixPath := flag.String("alterix", "", "Path to the alterix executable")
	inputDir := flag.String("input", "", "Input directory containing sigma rule files")
	outputDir := flag.String("output", "", "Output directory to store the processed JSON files")
	configPath := flag.String("config", "", "Path to the config file for the alterix command")

	// Parse flags
	flag.Parse()

	// Validate inputs
	if *alterixPath == "" || *inputDir == "" || *outputDir == "" || *configPath == "" {
		fmt.Println("Usage: -alterix <alterix-path> -input <input-dir> -output <output-dir> -config <config-path>")
		os.Exit(1)
	}

	// Process files and handle errors directly
	ProcessFile(*alterixPath, *inputDir, *outputDir, *configPath, "errors.log")

	fmt.Println("File processing complete.")
}
