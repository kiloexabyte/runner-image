package commands

import "log"

func (ops Ops) BuildAndUpload() {
	err := ops.Build()
	if err != nil {
		log.Fatal(err)
	}

	err = ops.Upload()
	if err != nil {
		log.Fatal(err)
	}
}
