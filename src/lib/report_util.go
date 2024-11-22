package lib

import (
	"fmt"
	"os"
	"strings"

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

func GenerateExcelFile(yamlFiles []string, outputExcel string) error {
	var tickets []Ticket

	for _, yamlFile := range yamlFiles {
		var tmpTickets []Ticket

		// Read the YAML file
		yamlFile, err := os.ReadFile(yamlFile)
		if err != nil {
			fmt.Printf("Error reading YAML file %s: %v\n", yamlFile, err)
			return err
		}

		// Parse the YAML data
		err = yaml.Unmarshal(yamlFile, &tmpTickets)
		if err != nil {
			fmt.Printf("Error unmarshaling YAML %s: %v\n", yamlFile, err)
			return err
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

	lastCol, err := drawHeader(f, sheetName, initialCol, currRow)
	if err != nil {
		fmt.Printf("Error drawing header: %s\n", err)
		return err
	}

	// Creating styles
	ticketHeaderStyleID, err := f.NewStyle(&excelize.Style{
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
		fmt.Printf("Error creating ticket header style: %s\n", err)
		return err
	}

	testcaseHeaderStyleID, err := f.NewStyle(&excelize.Style{
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
		fmt.Printf("Error creating testcase header style: %s\n", err)
		return err
	}

	leftCenterStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			WrapText:   true,
			Horizontal: "left",
			Vertical:   "center",
		},
		Border: allBorders,
	})

	if err != nil {
		fmt.Printf("Error creating testcase style: %s\n", err)
		return err
	}

	centerCenterStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			WrapText:   true,
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: allBorders,
	})

	if err != nil {
		fmt.Printf("Error creating centerCenter style: %s\n", err)
		return err
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
			f.SetCellValue(sheetName, getCell(currRow, startCol), "2024-01-04")

			// Tested By
			startCol = incrementColumnBy(endCol, 1)
			endCol = incrementColumnBy(startCol, TestedByColCount-1)
			f.SetCellStyle(sheetName, getCell(currRow, startCol), getCell(currRow, lastCol), centerCenterStyle)
			f.MergeCell(sheetName, getCell(currRow, startCol), getCell(currRow, endCol))
			f.SetCellValue(sheetName, getCell(currRow, startCol), "Paing Pyae Man")

		}
	}

	// Save the Excel file
	if err := f.SaveAs(outputExcel); err != nil {
		fmt.Printf("Error saving Excel file: %s\n", err)
		return err
	}

	fmt.Printf("Excel file successfully created as %s\n", outputExcel)

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
