package services

import (
	"bytes"
	"reflect"
	"slices"
	"strings"
	"text/template"

	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	servicetypes "github.com/otto-nation/otto-stack/internal/pkg/types"
)

// TemplateProcessor handles template processing for init scripts
type TemplateProcessor struct{}

// NewTemplateProcessor creates a new template processor
func NewTemplateProcessor() *TemplateProcessor {
	return &TemplateProcessor{}
}

// Process processes Go template variables in script content
func (tp *TemplateProcessor) Process(scriptContent string, config servicetypes.ServiceConfig, allConfigs []servicetypes.ServiceConfig) (string, error) {
	templateData := tp.collectTemplateData(config, allConfigs)

	tmpl, err := template.New("script").Parse(scriptContent)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, templateData); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (tp *TemplateProcessor) collectTemplateData(config servicetypes.ServiceConfig, allConfigs []servicetypes.ServiceConfig) map[string]any {
	templateData := make(map[string]any)

	for _, serviceConfig := range allConfigs {
		if tp.serviceDependsOn(serviceConfig, config.Name) {
			tp.addConfigData(templateData, serviceConfig)
		}
	}

	logger.GetLogger().Debug("Collected template data", "data", templateData)

	return templateData
}

func (tp *TemplateProcessor) serviceDependsOn(serviceConfig servicetypes.ServiceConfig, serviceName string) bool {
	if serviceConfig.Service.Dependencies.Required != nil {
		return slices.Contains(serviceConfig.Service.Dependencies.Required, serviceName)
	}
	return false
}

func (tp *TemplateProcessor) addConfigData(templateData map[string]any, serviceConfig servicetypes.ServiceConfig) {
	v := reflect.ValueOf(serviceConfig)

	for i := 0; i < v.NumField(); i++ {
		tp.processServiceField(templateData, v.Field(i))
	}
}

func (tp *TemplateProcessor) processServiceField(templateData map[string]any, field reflect.Value) {
	if field.Kind() != reflect.Pointer || field.IsNil() {
		return
	}

	structValue := field.Elem()
	structType := structValue.Type()

	if !strings.HasSuffix(structType.Name(), "Config") {
		return
	}

	tp.extractFieldsFromStruct(templateData, structValue, structType)
}

func (tp *TemplateProcessor) extractFieldsFromStruct(templateData map[string]any, structValue reflect.Value, structType reflect.Type) {
	for j := 0; j < structValue.NumField(); j++ {
		tp.processStructField(templateData, structValue.Field(j), structType.Field(j))
	}
}

func (tp *TemplateProcessor) processStructField(templateData map[string]any, structField reflect.Value, structFieldType reflect.StructField) {
	if !tp.isPopulatedSlice(structField) {
		return
	}

	fieldName := tp.getYAMLFieldName(structFieldType)
	if fieldName != "" {
		templateData[fieldName] = structField.Interface()
	}
}

func (tp *TemplateProcessor) isPopulatedSlice(field reflect.Value) bool {
	return field.Kind() == reflect.Slice && field.Len() > 0
}

func (tp *TemplateProcessor) getYAMLFieldName(structField reflect.StructField) string {
	yamlTag := structField.Tag.Get(ServiceCatalogYAMLFormat)
	fieldName := strings.Split(yamlTag, ",")[0]
	if fieldName != "-" {
		return fieldName
	}
	return ""
}
