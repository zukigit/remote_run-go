package common

type TestCase struct {
	Testcase_no          uint
	Testcase_description string
	Testcase_status      *Testcase_status
	Duration             *string
	Pre_operation        *[]string
	Operation            *[]string
	Expected_results     *[]string
	function             func() Testcase_status
	ticket_no            *uint
}

func New_testcase(testcase_id uint, testcase_description string) *TestCase {
	status := FAILED
	ticket_no := uint(0)
	pre_opt := []string{}
	opt := []string{}
	expt_res := []string{}
	var duration string

	return &TestCase{
		Testcase_no:          testcase_id,
		Testcase_description: testcase_description,
		Testcase_status:      &status,
		Pre_operation:        &pre_opt,
		Operation:            &opt,
		Expected_results:     &expt_res,
		ticket_no:            &ticket_no,
		Duration:             &duration,
	}
}

func (t *TestCase) Add_doc(doc_type Doc_data_type, doc string) {
	switch doc_type {
	case PRE_OPT:
		*t.Pre_operation = append(*t.Pre_operation, doc)
	case OPT:
		*t.Operation = append(*t.Operation, doc)
	case EXPT_RES:
		*t.Expected_results = append(*t.Expected_results, doc)
	}
}

func (t *TestCase) Get_testcase_no() uint {
	return t.Testcase_no
}

func (t *TestCase) Get_ticket_no() uint {
	return *t.ticket_no
}

func (t *TestCase) Set_ticket_no(ticket_no uint) {
	*t.ticket_no = ticket_no
}

func (t *TestCase) Get_ticket_description() string {
	return t.Testcase_description
}

func (t *TestCase) Set_status(status Testcase_status) {
	*t.Testcase_status = status
}

func (t *TestCase) Set_duration(duration string) {
	*t.Duration = duration
}

func (t *TestCase) Get_status() Testcase_status {
	return *t.Testcase_status
}

func (t *TestCase) Set_function(function func() Testcase_status) {
	t.function = function
}

func (t *TestCase) Run_function() Testcase_status {
	return t.function()
}

func (t *TestCase) Is_function_nil() bool {
	return t.function == nil
}
