package dao

type Jobnet_run_info struct {
	Jobnet_status, Job_status, Std_out, Std_error string
	Exit_cd                                       int64
}

func New_Jobnet(jobnet_status, job_status, std_out, std_error string, exit_cd int64) *Jobnet_run_info {
	return &Jobnet_run_info{
		Jobnet_status: jobnet_status,
		Job_status:    job_status,
		Exit_cd:       exit_cd,
		Std_out:       std_out,
		Std_error:     std_error,
	}
}
