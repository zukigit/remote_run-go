package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/zukigit/remote_run-go/src/lib"
)

const (
	logFileDir = "logs"
)

var testerName string

var generateExcelCommand = &cobra.Command{
	Use:   "generate-excel [yml files...]",
	Short: "Generate an Excel file from YAML files",
	Long:  `Parses the specified YAML files or all YAML files in the 'logs/' directory and generates an Excel file.`,
	Run: func(cmd *cobra.Command, args []string) {
		var yamlFiles []string
		var err error

		if len(args) == 0 {
			// no yaml files are specified, find and use all yml files under logs/
			logDir, err := getLogDir()
			if err != nil {
				log.Fatalf("failed to get logs directory: %v", err)
			}

			yamlFiles, err = findYAMLFiles(logDir)
			if err != nil {
				log.Fatalf("failed to find YAML files under %s directory: %v", logFileDir, err)
			}

			// check if yaml files are found
			if len(yamlFiles) == 0 {
				// return as no yaml files are found
				log.Fatalf("No YAML files found in %s directory.", logFileDir)
			}

		} else {
			yamlFiles = args
		}

		// convert the yaml files to an excel file
		reportDir, err := getReportDir()

		if err != nil {
			log.Fatalf("failed to get report directory: %v", err)
		}

		outputExcel := filepath.Join(reportDir, getExcelFileName())
		err = lib.GenerateExcelFile(yamlFiles, outputExcel, testerName)
		if err != nil {
			log.Fatalf("failed to generate excel: %v\n", err)
		}

		log.Printf("Excel file successfully created: %s\n", outputExcel)

	},
}

func getReportDir() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get exec path: %v", err)
	}

	execDir := filepath.Dir(execPath)

	reportDir := filepath.Join(execDir, "reports")

	err = os.MkdirAll(reportDir, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create directory %s: %v", reportDir, err)
	}

	return reportDir, nil
}

func getLogDir() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get exec path: %v", err)
	}

	execDir := filepath.Dir(execPath)

	logDir := filepath.Join(execDir, "logs")

	return logDir, nil
}

// this function an excel file name with the current timestamp
func getExcelFileName() string {
	// Get the current time
	currentTime := time.Now()

	// Format the time to "YYYYMMDDHHMMSS"
	timestamp := currentTime.Format("20060102150405")

	// Construct the file name
	fileName := fmt.Sprintf("%s_test_result.xlsx", timestamp)
	return fileName
}

func findYAMLFiles(rootPath string) ([]string, error) {
	var yamlFiles []string

	err := filepath.Walk(rootPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(info.Name()))
		if ext == ".yml" || ext == ".yaml" {
			yamlFiles = append(yamlFiles, path)
		}

		return nil
	})

	if err != nil {
		return yamlFiles, err
	}

	return yamlFiles, nil
}

func init() {
	generateExcelCommand.Flags().StringVarP(&testerName, "tester", "t", "", "Name of the tester")
	rootCmd.AddCommand(generateExcelCommand)
}
