package run

type RunOutput struct {
	Error error
}

func (o RunOutput) GetError() error {
	return o.Error
}
