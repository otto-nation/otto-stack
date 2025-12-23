//go:build unit

package unit

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/project"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/stack"
	"github.com/otto-nation/otto-stack/internal/pkg/config"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
)

// Performance benchmarks for critical operations

func BenchmarkStateManager_GetConfigHash(b *testing.B) {
	sm := stack.NewStateManager()
	cfg := &config.Config{
		Project: config.ProjectConfig{
			Name:     project.TestProjectName,
			Services: []string{services.ServicePostgres, services.ServiceRedis, services.ServiceMysql},
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

func BenchmarkFileGenerator_GenerateComposeFile(b *testing.B) {
	fg := services.NewFileGenerator()
	servicesList := []string{services.ServicePostgres, services.ServiceRedis, services.ServiceMysql}
	projectName := project.TestProjectName

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// This will fail because we can't write files, but we're measuring
		// the performance of the content generation logic
		_ = fg.GenerateComposeFile(servicesList, projectName)
	}
}

func BenchmarkFileGenerator_GenerateEnvFile(b *testing.B) {
	fg := services.NewFileGenerator()
	servicesList := []string{services.ServicePostgres, services.ServiceRedis, services.ServiceMysql}
	projectName := project.TestProjectName

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// This will fail because we can't write files, but we're measuring
		// the performance of the content generation logic
		_ = fg.GenerateEnvFile(servicesList, projectName)
	}
}

func BenchmarkStateManager_NewStateManager(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sm := stack.NewStateManager()
		_ = sm // Use sm to avoid unused variable
	}
}

// Memory allocation benchmarks
func BenchmarkStateManager_GetConfigHash_Memory(b *testing.B) {
	sm := stack.NewStateManager()
	cfg := &config.Config{
		Project: config.ProjectConfig{
			Name:     "memory-benchmark-project",
			Services: make([]string, 100), // Large service list
		},
	}

	// Fill with service names
	for i := 0; i < 100; i++ {
		cfg.Project.Services[i] = "service-" + string(rune(i))
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
