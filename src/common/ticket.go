package common

type Ticket interface {
	Get_ticket_no() uint
	Set_ticket_no(ticket_no uint)
	Get_ticket_description() string
	Set_ticket_description(testcase_description string)
	Set_PASSED_count(passed_count int)
	Set_FAILED_count(failed_count int)
	Set_MUSTCHECK_count(mustcheck_count int)
	Add_testcase(tc *TestCase)
	Prepare()
	Get_testcases() []TestCase
	New_testcase(testcase_id uint, testcase_description string) *TestCase
}

type TicketStruct struct {
	TicketNo          int        `yaml:"ticket_no"`
	TicketDescription string     `yaml:"ticket_description"`
	PassedCount       int        `yaml:"passed_count"`
	FailedCount       int        `yaml:"failed_count"`
	MustCheckCount    int        `yaml:"mustcheck_count"`
	Testcases         []TestCase `yaml:"testcases"`
	TestedDate        string
}
