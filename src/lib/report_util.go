package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
	"gopkg.in/yaml.v3"
)

// number of columns to be merged
const (
	NoColCount              = 3
	PreOperationColCount    = 17
	OperationColCount       = 17
	ExpectedResultsColCount = 17
	DurationColCount        = 2
	ResultColCount          = 1
	TestedDateColCount      = 3
	TestedByColCount        = 3
)

// Define the structure to parse the YAML
type TestCase struct {
	TestcaseNo          int      `yaml:"testcase_no"`
	TestcaseDescription string   `yaml:"testcase_description"`
	PreOperation        []string `yaml:"pre_operation"`
	Operation           []string `yaml:"operation"`
	ExpectedResults     []string `yaml:"expected_results"`
	TestcaseStatus      string   `yaml:"testcase_status"`
}

type Ticket struct {
	TicketNo          int        `yaml:"ticket_no"`
	TicketDescription string     `yaml:"ticket_description"`
	PassedCount       int        `yaml:"passed_count"`
	FailedCount       int        `yaml:"failed_count"`
	MustCheckCount    int        `yaml:"mustcheck_count"`
	Testcases         []TestCase `yaml:"testcases"`
	TestedDate        string
}

func (t *TestCase) getResult() string {
	if t.TestcaseStatus == "PASSED" {
		return "O"
	} else if t.TestcaseStatus == "FAILED" {
		return "X"
	} else if t.TestcaseStatus == "MUST_CHECK" {
		return "-"
	} else {
		return "-"
	}
}

var allBorders []excelize.Border
var ticketHeaderStyleID, testcaseHeaderStyleID, leftCenterStyle, centerCenterStyle int

func init() {
	allBorders = []excelize.Border{
		{
			Type:  "left",
			Color: "000000", // Black color
			Style: 1,        // Thin border
		},
		{
			Type:  "right",
			Color: "000000",
			Style: 1, // Thin border
		},
		{
			Type:  "top",
			Color: "000000",
			Style: 1, // Thin border
		},
		{
			Type:  "bottom",
			Color: "000000",
			Style: 1, // Thin border
		},
	}
}

func GenerateExcelFile(yamlFiles []string, outputExcel, testerName string) error {
	var tickets []Ticket

	for _, yamlFilePath := range yamlFiles {
		var tmpTickets []Ticket

		// Read the YAML file
		yamlFile, err := os.ReadFile(yamlFilePath)
		if err != nil {
			return fmt.Errorf("error reading YAML file %s: %v", yamlFile, err)
		}

		// Parse the YAML data
		err = yaml.Unmarshal(yamlFile, &tmpTickets)
		if err != nil {
			return fmt.Errorf("error unmarshaling YAML %s: %v", yamlFile, err)
		}

		// get Tested Date
		testedDate, err := getTestedDate(yamlFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to set 'testedDate': %v", err)
		}

		for i := range tmpTickets {
			tmpTickets[i].TestedDate = testedDate
		}

		tickets = append(tickets, tmpTickets...)
	}

	// Create a new Excel file
	f := excelize.NewFile()

	// Build Up headers
	sheetName := "論理テスト"
	f.NewSheet(sheetName)
	currRow := 2
	initialCol := "B"

	// Creating styles
	if err := initializeStyles(f); err != nil {
		return fmt.Errorf("failed to initialize styles: %v", err)
	}

	// Drawing header
	lastCol, err := drawHeader(f, sheetName, initialCol, currRow)
	if err != nil {
		return fmt.Errorf("error drawing header: %v", err)
	}

	// Create a sheet for the ticket information
	for _, ticket := range tickets {
		fmt.Printf("Generating for ticket [%d] %s\n", ticket.TicketNo, ticket.TicketDescription)

		// Write ticket no
		currRow++
		startCol := initialCol
		f.SetCellValue(sheetName, getCell(currRow, startCol), ticket.TicketNo)

		// Write ticket description
		startCol = incrementColumnBy(startCol, 1)
		f.MergeCell(sheetName, getCell(currRow, startCol), getCell(currRow, lastCol))
		f.SetCellValue(sheetName, getCell(currRow, startCol), ticket.TicketDescription)
		f.SetCellStyle(sheetName, getCell(currRow, initialCol), getCell(currRow, lastCol), ticketHeaderStyleID)

		// Populate the test cases
		for _, testcase := range ticket.Testcases {
			fmt.Printf("    Generating for testcase [%d] %s\n", testcase.TestcaseNo, testcase.TestcaseDescription)

			// Write testcase no
			currRow++
			startCol = incrementColumnBy(initialCol, 1)
			f.SetCellValue(sheetName, getCell(currRow, startCol), testcase.TestcaseNo)

			// Write testcase description
			startCol = incrementColumnBy(startCol, 1)
			f.MergeCell(sheetName, getCell(currRow, startCol), getCell(currRow, lastCol))
			f.SetCellValue(sheetName, getCell(currRow, startCol), testcase.TestcaseDescription)

			f.SetCellStyle(sheetName, getCell(currRow, initialCol), getCell(currRow, lastCol), testcaseHeaderStyleID)

			// set a suitable height according to number of lines in a single cell
			maxLines := getMaxNumber(
				len(testcase.PreOperation),
				len(testcase.Operation),
				len(testcase.ExpectedResults),
			)

			if maxLines <= 0 {
				continue
			}

			// Write testcase details, such as pre-operations, operations, and expected results
			currRow++

			f.SetRowHeight(sheetName, currRow, 1.5*11*float64(maxLines))

			// Write pre-operation
			startCol = incrementColumnBy(startCol, 1)
			endCol := incrementColumnBy(startCol, PreOperationColCount-1)
			f.MergeCell(sheetName, getCell(currRow, startCol), getCell(currRow, endCol))
			f.SetCellValue(sheetName, getCell(currRow, startCol), strings.Join(testcase.PreOperation, "\n"))

			// Write operation
			startCol = incrementColumnBy(endCol, 1)
			endCol = incrementColumnBy(startCol, OperationColCount-1)
			f.MergeCell(sheetName, getCell(currRow, startCol), getCell(currRow, endCol))
			f.SetCellValue(sheetName, getCell(currRow, startCol), strings.Join(testcase.Operation, "\n"))

			// Write expected results
			startCol = incrementColumnBy(endCol, 1)
			endCol = incrementColumnBy(startCol, OperationColCount-1)
			f.MergeCell(sheetName, getCell(currRow, startCol), getCell(currRow, endCol))
			f.SetCellValue(sheetName, getCell(currRow, startCol), strings.Join(testcase.ExpectedResults, "\n"))

			f.SetCellStyle(sheetName, getCell(currRow, initialCol), getCell(currRow, endCol), leftCenterStyle)

			// Write expected duration
			startCol = incrementColumnBy(endCol, 1)
			endCol = incrementColumnBy(startCol, DurationColCount-1)
			f.SetCellStyle(sheetName, getCell(currRow, startCol), getCell(currRow, lastCol), centerCenterStyle)
			f.MergeCell(sheetName, getCell(currRow, startCol), getCell(currRow, endCol))
			f.SetCellValue(sheetName, getCell(currRow, startCol), 15)

			// Write actual duration
			startCol = incrementColumnBy(endCol, 1)
			endCol = incrementColumnBy(startCol, DurationColCount-1)
			f.SetCellStyle(sheetName, getCell(currRow, startCol), getCell(currRow, lastCol), centerCenterStyle)
			f.MergeCell(sheetName, getCell(currRow, startCol), getCell(currRow, endCol))
			f.SetCellValue(sheetName, getCell(currRow, startCol), 20)

			// Write result
			startCol = incrementColumnBy(endCol, 1)
			endCol = incrementColumnBy(startCol, ResultColCount-1)
			f.SetCellStyle(sheetName, getCell(currRow, startCol), getCell(currRow, lastCol), centerCenterStyle)
			f.MergeCell(sheetName, getCell(currRow, startCol), getCell(currRow, endCol))
			f.SetCellValue(sheetName, getCell(currRow, startCol), testcase.getResult())

			// Tested Date
			startCol = incrementColumnBy(endCol, 1)
			endCol = incrementColumnBy(startCol, TestedDateColCount-1)
			f.SetCellStyle(sheetName, getCell(currRow, startCol), getCell(currRow, lastCol), centerCenterStyle)
			f.MergeCell(sheetName, getCell(currRow, startCol), getCell(currRow, endCol))
			f.SetCellValue(sheetName, getCell(currRow, startCol), ticket.TestedDate)

			// Tested By
			startCol = incrementColumnBy(endCol, 1)
			endCol = incrementColumnBy(startCol, TestedByColCount-1)
			f.SetCellStyle(sheetName, getCell(currRow, startCol), getCell(currRow, lastCol), centerCenterStyle)
			f.MergeCell(sheetName, getCell(currRow, startCol), getCell(currRow, endCol))
			f.SetCellValue(sheetName, getCell(currRow, startCol), testerName)

		}

	}

	f.SetColWidth(sheetName, initialCol, initialCol, 5)

	// Save the Excel file
	if err := f.SaveAs(outputExcel); err != nil {
		return fmt.Errorf("error saving Excel file: %s", err)
	}

	return nil
}

func initializeStyles(f *excelize.File) error {
	var err error
	ticketHeaderStyleID, err = f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			WrapText:   true,
			Horizontal: "left",
			Vertical:   "center",
		},
		Border: allBorders,
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#C5D9F1"},
			Pattern: 1,
		},
	})

	if err != nil {
		return fmt.Errorf("error creating ticket header style: %w", err)
	}

	testcaseHeaderStyleID, err = f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			WrapText:   true,
			Horizontal: "left",
			Vertical:   "center",
		},
		Border: allBorders,
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#F2DCDB"},
			Pattern: 1,
		},
	})

	if err != nil {
		return fmt.Errorf("error creating testcase header style: %w", err)
	}

	leftCenterStyle, err = f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			WrapText:   true,
			Horizontal: "left",
			Vertical:   "center",
		},
		Border: allBorders,
	})

	if err != nil {
		return fmt.Errorf("error creating testcase style: %w", err)
	}

	centerCenterStyle, err = f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			WrapText:   true,
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: allBorders,
	})

	if err != nil {
		return fmt.Errorf("error creating centerCenter style: %w", err)
	}

	return nil
}

func getCell(row int, col string) string {
	return fmt.Sprintf("%s%d", col, row)
}

func drawHeader(f *excelize.File, sheetName, initialCol string, headerRow int) (string, error) {
	f.NewSheet(sheetName)

	startCol := initialCol

	endCol := incrementColumnBy(startCol, NoColCount-1)

	f.MergeCell(
		sheetName,
		getCell(headerRow, startCol),
		getCell(headerRow, endCol),
	)
	f.SetCellValue(sheetName, getCell(headerRow, startCol), "No.")

	startCol = incrementColumnBy(endCol, 1)
	endCol = incrementColumnBy(startCol, PreOperationColCount-1)
	f.MergeCell(
		sheetName,
		getCell(headerRow, startCol),
		getCell(headerRow, endCol),
	)
	f.SetCellValue(sheetName, getCell(headerRow, startCol), "Pre-Operation State\n事前作業状態")

	startCol = incrementColumnBy(endCol, 1)
	endCol = incrementColumnBy(startCol, OperationColCount-1)
	f.MergeCell(
		sheetName,
		getCell(headerRow, startCol),
		getCell(headerRow, endCol),
	)
	f.SetCellValue(sheetName, getCell(headerRow, startCol), "Operation\n作業")

	startCol = incrementColumnBy(endCol, 1)
	endCol = incrementColumnBy(startCol, ExpectedResultsColCount-1)
	f.MergeCell(
		sheetName,
		getCell(headerRow, startCol),
		getCell(headerRow, endCol),
	)
	f.SetCellValue(sheetName, getCell(headerRow, startCol), "Expected Results\n期待結果")

	startCol = incrementColumnBy(endCol, 1)
	endCol = incrementColumnBy(startCol, DurationColCount-1)
	f.MergeCell(
		sheetName,
		getCell(headerRow, startCol),
		getCell(headerRow, endCol),
	)
	f.SetCellValue(sheetName, getCell(headerRow, startCol), "Expected Duration\n(min)")

	startCol = incrementColumnBy(endCol, 1)
	endCol = incrementColumnBy(startCol, DurationColCount-1)
	f.MergeCell(
		sheetName,
		getCell(headerRow, startCol),
		getCell(headerRow, endCol),
	)
	f.SetCellValue(sheetName, getCell(headerRow, startCol), "Actual Duration\n(min)")

	startCol = incrementColumnBy(endCol, 1)
	endCol = incrementColumnBy(startCol, ResultColCount-1)
	f.SetCellValue(sheetName, getCell(headerRow, startCol), "O/X")

	startCol = incrementColumnBy(endCol, 1)
	endCol = incrementColumnBy(startCol, TestedDateColCount-1)
	f.MergeCell(
		sheetName,
		getCell(headerRow, startCol),
		getCell(headerRow, endCol),
	)
	f.SetCellValue(sheetName, getCell(headerRow, startCol), "Tested Date\nテスト日")

	startCol = incrementColumnBy(endCol, 1)
	endCol = incrementColumnBy(startCol, TestedByColCount-1)
	f.MergeCell(
		sheetName,
		getCell(headerRow, startCol),
		getCell(headerRow, endCol),
	)
	f.SetCellValue(sheetName, getCell(headerRow, startCol), "Tested By\nテスト者")

	style, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			WrapText:   true,
			Horizontal: "center",
			Vertical:   "center",
		},
		Font: &excelize.Font{
			Size:  10,
			Bold:  true,
			Color: "#FFFFFF",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"365F92"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{
				Type:  "left",
				Color: "000000", // Black color
				Style: 1,        // Thin border
			},
			{
				Type:  "right",
				Color: "000000",
				Style: 1, // Thin border
			},
			{
				Type:  "top",
				Color: "000000",
				Style: 1, // Thin border
			},
			{
				Type:  "bottom",
				Color: "000000",
				Style: 1, // Thin border
			},
		},
	})

	if err != nil {
		return "", fmt.Errorf("failed to create style: %w", err)
	}
	f.SetCellStyle(sheetName, getCell(headerRow, initialCol), getCell(headerRow, endCol), style)

	f.SetColWidth(sheetName, "A", "ZZ", 4)
	f.SetRowHeight(sheetName, headerRow, 60)

	finalCol := endCol

	return finalCol, nil

}

// Helper function to convert column label to number
func columnToNumber(col string) int {
	col = strings.ToUpper(col)
	num := 0
	for _, char := range col {
		num = num*26 + int(char-'A'+1)
	}
	return num
}

// Helper function to convert number back to column label
func numberToColumn(num int) string {
	var col string
	for num > 0 {
		num-- // Adjust for 1-based indexing
		col = string(rune(num%26+'A')) + col
		num /= 26
	}
	return col
}

// Increment function to add a specified number of steps
func incrementColumnBy(col string, steps int) string {
	currentNumber := columnToNumber(col)
	newNumber := currentNumber + steps
	return numberToColumn(newNumber)
}

// get the maximum number
func getMaxNumber(numbers ...int) int {
	var maxNumber int

	for i, number := range numbers {
		if i == 0 || number > maxNumber {
			maxNumber = number
		}
	}

	return maxNumber
}

// get test date 'YYYY-MM-dd' from yaml file name
func getTestedDate(yamlFilePath string) (string, error) {
	yamlFileName := filepath.Base(yamlFilePath)
	parts := strings.Split(yamlFileName, "_")

	if len(parts) <= 0 {
		return "", fmt.Errorf("yaml file name does not contain '_'")
	}

	// parsing the timestamp
	timestamp := parts[0]
	inputLayout := "20060102150405.000"

	parsedTime, err := time.Parse(inputLayout, timestamp)
	if err != nil {
		return "", fmt.Errorf("failed to parse timestamp: %w", err)
	}

	// convert the output format
	outputLayout := "2006-01-02"
	output := parsedTime.Format(outputLayout)
	return output, nil
}
