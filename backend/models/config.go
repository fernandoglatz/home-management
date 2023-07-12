package models

type Config struct {
	Database struct {
		URI  string `yaml:"uri"`
		Name string `yaml:"name"`
	} `yaml:"database"`
	Server struct {
		Listening string `yaml:"listening"`
	} `yaml:"server"`
}
