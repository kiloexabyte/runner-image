package commands

func (ops Ops) BuildAndUpload() error {
	err := ops.Build()
	if err != nil {
		return err
	}

	err = ops.Upload()
	if err != nil {
		return err
	}

	return nil
}
