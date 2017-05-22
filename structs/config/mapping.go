package config

type StaticEntry struct {
	ServerUrl *string `mapstructure:"server"`
}

type Entry struct {
	Name        *string `mapstructure:`
	ServiceName *string `mapstructure:"service"`
	Static      []*StaticEntry
}

type Mapping struct {
	Entries []*Entry
}
