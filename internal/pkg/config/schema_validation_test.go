package config

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestSharingConfig_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		wantErr bool
	}{
		{
			name: "valid sharing config",
			yaml: `
sharing:
  enabled: true
  services:
    postgres: true
    redis: false
`,
			wantErr: false,
		},
		{
			name: "sharing disabled",
			yaml: `
sharing:
  enabled: false
`,
			wantErr: false,
		},
		{
			name: "invalid enabled type",
			yaml: `
sharing:
  enabled: "not-a-bool"
`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cfg struct {
				Sharing *SharingConfig `yaml:"sharing"`
			}
			err := yaml.Unmarshal([]byte(tt.yaml), &cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
