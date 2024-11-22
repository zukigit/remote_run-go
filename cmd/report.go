package cmd

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zukigit/remote_run-go/src/lib"
)

const (
	logFileDir = "logs"
)

var generateExcelCommand = &cobra.Command{
	Use:   "generate-excel [yml files...]",
	Short: "Generate an Excel file from YAML files",
	Long:  `Parses the specified YAML files or all YAML files in the 'log/' directory and generates an Excel file.`,
	Run: func(cmd *cobra.Command, args []string) {
		var yamlFiles []string
		var err error

		if len(args) == 0 {
			// no yaml files are specified, find and use all yml files under log/
			yamlFiles, err = findYAMLFiles(logFileDir)
			if err != nil {
				fmt.Printf("failed to find YAML files under %s directory.", logFileDir)
			}

			// check if yaml files are found
			if len(yamlFiles) == 0 {
				// return as no yaml files are found
				fmt.Printf("No YAML files found in %s directory.", logFileDir)
				return
			}

		} else {
			yamlFiles = args
		}

		// convert the yaml files to an excel file
		outputExcel := "logs/test-output.xlsx"
		err = lib.GenerateExcelFile(yamlFiles, outputExcel)
		if err != nil {
			fmt.Printf("failed to generate excel: %v\n", err)
		}

	},
}

func findYAMLFiles(rootPath string) ([]string, error) {
	var yamlFiles []string

	err := filepath.Walk(rootPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
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
	rootCmd.AddCommand(generateExcelCommand)
}
