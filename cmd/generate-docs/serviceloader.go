package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type serviceConfig struct {
	Name               string                `yaml:"name"`
	Description        string                `yaml:"description"`
	Hidden             bool                  `yaml:"hidden"`
	Environment        map[string]string     `yaml:"environment"`
	Documentation      *serviceDocumentation `yaml:"documentation"`
	configSchemaFields []*schemaField
}

type serviceDocumentation struct {
	UseCases []string `yaml:"use_cases"`
	Examples []string `yaml:"examples"`
}

type loadedService struct {
	name     string
	config   serviceConfig
	category string
}

type categoryConfig struct {
	Icon  string `yaml:"icon"`
	Order int    `yaml:"order"`
}

func getCategoryConfig(name string) categoryConfig {
	if c, ok := docs.Categories[name]; ok {
		return c
	}
	return docs.Categories["other"]
}

func loadAllServices() ([]loadedService, error) {
	var services []loadedService
	err := filepath.Walk(servicesDirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		if ext != ".yaml" && ext != ".yml" {
			return nil
		}
		svc, err := loadService(path)
		if err != nil {
			return err
		}
		if svc != nil {
			services = append(services, *svc)
		}
		return nil
	})
	return services, err
}

func loadService(path string) (*loadedService, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}

	var rootNode yaml.Node
	if err := yaml.Unmarshal(data, &rootNode); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}

	var svc serviceConfig
	if err := nodeDoc(&rootNode).Decode(&svc); err != nil {
		return nil, fmt.Errorf("decode %s: %w", path, err)
	}
	if svc.Hidden {
		return nil, nil
	}

	svc.configSchemaFields = extractSchemaFields(nodeGet(&rootNode, keyConfigSchema))

	ext := filepath.Ext(path)
	name := strings.TrimSuffix(filepath.Base(path), ext)
	return &loadedService{name: name, config: svc, category: inferCategory(path)}, nil
}

func inferCategory(path string) string {
	relPath, _ := filepath.Rel(servicesDirPath, path)
	parts := strings.Split(filepath.ToSlash(relPath), "/")
	if len(parts) >= 2 {
		return parts[0]
	}
	return "other"
}

func indexServices(services []loadedService) map[string]loadedService {
	svcMap := make(map[string]loadedService, len(services))
	for _, svc := range services {
		svcMap[svc.name] = svc
	}
	return svcMap
}

func groupByCategory(services []loadedService) map[string][]loadedService {
	byCategory := make(map[string][]loadedService)
	for _, svc := range services {
		byCategory[svc.category] = append(byCategory[svc.category], svc)
	}
	return byCategory
}

func sortedCategories(byCategory map[string][]loadedService) []string {
	categories := make([]string, 0, len(byCategory))
	for cat := range byCategory {
		categories = append(categories, cat)
	}
	sort.Slice(categories, func(i, j int) bool {
		return getCategoryConfig(categories[i]).Order < getCategoryConfig(categories[j]).Order
	})
	return categories
}
