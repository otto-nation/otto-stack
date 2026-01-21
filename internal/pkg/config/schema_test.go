package config

import "testing"

func TestSharingConfig_Defaults(t *testing.T) {
	cfg := &SharingConfig{
		Enabled: true,
	}

	if !cfg.Enabled {
		t.Error("expected Enabled to be true")
	}
	if cfg.Services != nil {
		t.Error("expected Services to be nil by default")
	}
}

func TestSharingConfig_WithServices(t *testing.T) {
	cfg := &SharingConfig{
		Enabled: true,
		Services: map[string]bool{
			"postgres": true,
			"redis":    false,
		},
	}

	if !cfg.Enabled {
		t.Error("expected Enabled to be true")
	}
	if cfg.Services["postgres"] != true {
		t.Error("expected postgres to be true")
	}
	if cfg.Services["redis"] != false {
		t.Error("expected redis to be false")
	}
}

func TestConfig_WithSharing(t *testing.T) {
	cfg := &Config{
		Project: ProjectConfig{
			Name: "test",
			Type: "development",
		},
		Stack: StackConfig{
			Enabled: []string{"postgres", "redis"},
		},
		Sharing: &SharingConfig{
			Enabled: true,
			Services: map[string]bool{
				"postgres": true,
			},
		},
	}

	if cfg.Sharing == nil {
		t.Fatal("expected Sharing to be populated")
	}
	if !cfg.Sharing.Enabled {
		t.Error("expected Sharing.Enabled to be true")
	}
	if cfg.Sharing.Services["postgres"] != true {
		t.Error("expected postgres sharing to be true")
	}
}
