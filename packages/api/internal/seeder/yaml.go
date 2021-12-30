package seeder

import "gopkg.in/yaml.v3"

type Parsed struct {
	Phones map[string]string `yaml:"phones"`
}

func ParseYamlFromString(yamlData string) (Parsed, error) {
	out := Parsed{}

	err := yaml.Unmarshal([]byte(yamlData), &out)
	return out, err
}
