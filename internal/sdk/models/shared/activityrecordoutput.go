// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package shared

// Job - State of the job to post-process the records.
type Job struct {
	// The ID of the job.
	JobID string `json:"job_id"`
	// The state of the job.
	State *RunInfoStatus `json:"state,omitempty"`
}

func (o *Job) GetJobID() string {
	if o == nil {
		return ""
	}
	return o.JobID
}

func (o *Job) GetState() *RunInfoStatus {
	if o == nil {
		return nil
	}
	return o.State
}

type ActivityRecordOutput struct {
	// State of the job to post-process the records.
	Job Job `json:"job"`
}

func (o *ActivityRecordOutput) GetJob() Job {
	if o == nil {
		return Job{}
	}
	return o.Job
}
