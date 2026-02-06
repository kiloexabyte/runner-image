package commands

import "fmt"

func (ops Ops) BuildAndUpload() error {
	err := ops.Build()
	if err != nil {
		return fmt.Errorf("build: %w", err)
	}

	err = ops.Upload()
	if err != nil {
		return fmt.Errorf("upload: %w", err)
	}

	return nil
}
