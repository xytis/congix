package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/containous/flaeg"
	"os"
)

type CongixConfiguration struct {
	Name string
}

func NewCongixConfiguration() *CongixConfiguration {
	return &CongixConfiguration{}
}

func NewPointersCongixConfiguration() *CongixConfiguration {
	return &CongixConfiguration{}
}

func main() {
	congixConfiguration := NewCongixConfiguration()
	pointersCongixConfiguration := NewPointersCongixConfiguration()
	congixCmd := &flaeg.Command{
		Name:                  "congix",
		Description:           `nginx-plus adapter to consul`,
		Config:                congixConfiguration,
		DefaultPointersConfig: pointersCongixConfiguration,
		Run: func() error {
			log.Error("Hello world")
			return nil
		},
	}

	flaeg := flaeg.New(congixCmd, os.Args[1:])
	if err := flaeg.Run(); err != nil {
		log.Fatalf("Error %s", err.Error())
	}
}
