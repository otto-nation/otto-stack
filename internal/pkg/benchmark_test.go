//go:build unit

package pkg

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/project"
	"github.com/otto-nation/otto-stack/internal/pkg/config"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
)

// Performance benchmarks for critical operations

func BenchmarkStateManager_GetConfigHash(b *testing.B) {
	sm := common.NewStateManager()
	cfg := &config.Config{
		Project: config.ProjectConfig{
			Name: project.TestProjectName,
		},
		Stack: config.StackConfig{
			Enabled: []string{services.ServicePostgres, services.ServiceRedis, services.ServiceMysql},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := sm.GetConfigHash(cfg)
		if err != nil {
			b.Fatalf("GetConfigHash failed: %v", err)
		}
	}
}

func BenchmarkStateManager_NewStateManager(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sm := common.NewStateManager()
		_ = sm // Use sm to avoid unused variable
	}
}

// Memory allocation benchmarks
func BenchmarkStateManager_GetConfigHash_Memory(b *testing.B) {
	sm := common.NewStateManager()
	cfg := &config.Config{
		Project: config.ProjectConfig{
			Name: "memory-benchmark-project",
		},
		Stack: config.StackConfig{
			Enabled: make([]string, 100), // Large service list
		},
	}

	// Fill with service names
	for i := 0; i < 100; i++ {
		cfg.Stack.Enabled[i] = "service-" + string(rune(i))
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := sm.GetConfigHash(cfg)
		if err != nil {
			b.Fatalf("GetConfigHash failed: %v", err)
		}
	}
}
