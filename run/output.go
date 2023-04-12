package run

import "fmt"

type RunOutput struct {
	error          error
	errorMessage   string
	successMessage string
}

func NewRunOutput(error error, errorMessage string, successMessage string) RunOutput {
	return RunOutput{
		error:          error,
		errorMessage:   errorMessage,
		successMessage: successMessage,
	}
}

func (o RunOutput) HasError() bool {
	return o.error != nil
}

func (o RunOutput) GetErrorMessage() string {
	return fmt.Sprintf("%s: %s", o.errorMessage, o.error)
}

func (o RunOutput) GetSuccessMessage() string {
	return o.successMessage
}
