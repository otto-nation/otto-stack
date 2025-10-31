package display

import (
	"fmt"
	"io"
	"strings"

	"github.com/otto-nation/otto-stack/internal/pkg/constants"
)

// CreateFormatter creates a formatter based on the specified format
func CreateFormatter(format string, writer io.Writer) (Formatter, error) {
	switch strings.ToLower(format) {
	case constants.ServiceCatalogTableFormat:
		return NewTableFormatter(writer), nil
	case constants.ServiceCatalogGroupFormat, "":
		return NewGroupFormatter(writer), nil
	case constants.ServiceCatalogJSONFormat:
		return NewJSONFormatter(writer), nil
	case constants.ServiceCatalogYAMLFormat, "yml":
		return NewYAMLFormatter(writer), nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// GetSupportedFormats returns a list of supported output formats
func GetSupportedFormats() []string {
	return []string{
		constants.ServiceCatalogGroupFormat,
		constants.ServiceCatalogTableFormat,
		constants.ServiceCatalogJSONFormat,
		constants.ServiceCatalogYAMLFormat,
	}
}
