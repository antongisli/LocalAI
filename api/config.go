package api

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	OpenAIRequest
	Name       string   `yaml:"name"`
	Path       string   `yaml:"path"`
	StopWords  []string `yaml:"stopwords"`
	Cutstrings []string `yaml:"cutstrings"`
	TrimSpace  []string `yaml:"trimspace"`

	ContextSize    int               `yaml:"context_size"`
	F16            bool              `yaml:"f16"`
	Threads        int               `yaml:"threads"`
	Debug          bool              `yaml:"debug"`
	Roles          map[string]string `yaml:"roles"`
	TemplateConfig TemplateConfig    `yaml:"template"`
}

type TemplateConfig struct {
	Completion string `yaml:"completion"`
	Chat       string `yaml:"chat"`
}

type ConfigFile struct {
	Configs []*Config
}

type ConfigMerger map[string]Config

func ReadConfigFile(file string) (*ConfigFile, error) {
	c := &ConfigFile{}
	f, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("cannot read config file: %w", err)
	}
	if err := yaml.Unmarshal(f, c); err != nil {
		return nil, fmt.Errorf("cannot unmarshal config file: %w", err)
	}

	return c, nil
}

func ReadConfig(file string) (*Config, error) {
	c := &Config{}
	f, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("cannot read config file: %w", err)
	}
	if err := yaml.Unmarshal(f, c); err != nil {
		return nil, fmt.Errorf("cannot unmarshal config file: %w", err)
	}

	return c, nil
}

func (cm ConfigMerger) LoadConfigFile(file string) error {
	c, err := ReadConfig(file)
	if err != nil {
		return fmt.Errorf("cannot read config file: %w", err)
	}

	cm[c.Name] = *c
	return nil
}

func (cm ConfigMerger) LoadConfig(file string) error {
	c, err := ReadConfigFile(file)
	if err != nil {
		return fmt.Errorf("cannot read config file: %w", err)
	}

	for _, cc := range c.Configs {
		cm[cc.Name] = *cc
	}
	return nil
}

func (cm ConfigMerger) LoadConfigs(path string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		// Skip templates, YAML and .keep files
		if !strings.Contains(file.Name(), ".yaml") {
			continue
		}
		c, err := ReadConfig(filepath.Join(path, file.Name()))
		if err == nil {
			cm[c.Name] = *c
		}
	}

	return nil
}
