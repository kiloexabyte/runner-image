package commands

import "log"

func (Ops) BuildAndUpload() {
	err := Ops{}.Build()
	if err != nil {
		log.Fatal(err)
	}

	err = Ops{}.Upload()
	if err != nil {
		log.Fatal(err)
	}
}
