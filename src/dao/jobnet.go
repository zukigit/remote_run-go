package dao

type Jobnet struct {
	Status, Job_status, Exit_cd, Std_out, Std_error string
}

func New_Jobnet(status, job_status, exit_cd, std_out, std_error string) *Jobnet {
	return &Jobnet{
		Status:     status,
		Job_status: job_status,
		Exit_cd:    exit_cd,
		Std_out:    std_out,
		Std_error:  std_error,
	}
}
