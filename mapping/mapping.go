package mapping

import (
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
)

type Mapping struct {
	Entries map[string]*Entry
}

type Entry interface {
	entry()
}

func (StaticEntry) entry() {}
func (ConsulEntry) entry() {}

type StaticEntry struct {
	Name      string
	Upstreams []string
}

func NewStaticEntry(name string, source map[string]interface{}) (*StaticEntry, error) {
	var entry StaticEntry
	if err := mapstructure.WeakDecode(source, &entry); err != nil {
		return nil, err
	}
	entry.Name = name
	return &entry, nil
}

type ConsulEntry struct {
	Name        string
	Service     string
	DelayRemove time.Duration `mapstructure:"delay_remove"`
	DelayInsert time.Duration `mapstructure:"delay_insert"`
}

func NewConsulEntry(name string, source map[string]interface{}) (*ConsulEntry, error) {
	var entry ConsulEntry
	config := &mapstructure.DecoderConfig{
		DecodeHook: mapstructure.StringToTimeDurationHookFunc(),
		Metadata:   nil,
		Result:     &entry,
		//WeaklyTypedInput: true,
	}
	fmt.Printf("---\nbefore: %v\n", source)
	if d, err := mapstructure.NewDecoder(config); err != nil {
		return nil, err
	} else {
		if err := d.Decode(source); err != nil {
			return nil, err
		}
		entry.Name = name
		fmt.Printf("after: %v\n---\n", entry)
		return &entry, nil
	}
}

func NewEntry(name string, source map[string]interface{}) (Entry, error) {
	if _, ok := source["type"]; !ok {
		return nil, fmt.Errorf("'entry' must have 'type' defined")
	}
	switch t := source["type"]; t {
	case "static":
		return NewStaticEntry(name, source)
	case "consul":
		return NewConsulEntry(name, source)
	default:
		return nil, fmt.Errorf("unsuported type value '%v'", t)
	}
}
