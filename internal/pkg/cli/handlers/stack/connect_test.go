package stack

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/services"
)

func TestConnectHandler_ValidateArgs(t *testing.T) {
	handler := NewConnectHandler()

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "valid args with service name",
			args:    []string{services.ServicePostgres},
			wantErr: false,
		},
		{
			name:    "missing service name",
			args:    []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.ValidateArgs(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateArgs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConnectHandler_GetRequiredFlags(t *testing.T) {
	handler := NewConnectHandler()
	flags := handler.GetRequiredFlags()

	if len(flags) != 0 {
		t.Errorf("GetRequiredFlags() = %v, want empty slice", flags)
	}
}

func TestConnectHandler_getConnectionCommand(t *testing.T) {
	handler := NewConnectHandler()

	tests := []struct {
		name        string
		serviceName string
		database    string
		user        string
		host        string
		port        int
		readOnly    bool
		wantCmd     []string
		wantErr     bool
	}{
		{
			name:        "postgres with defaults",
			serviceName: services.ServicePostgres,
			wantCmd:     []string{services.ClientPsql, services.User_flagPostgres, services.DefaultUserPostgres},
			wantErr:     false,
		},
		{
			name:        "postgres with custom user",
			serviceName: services.ServicePostgres,
			user:        "custom",
			wantCmd:     []string{services.ClientPsql, services.User_flagPostgres, "custom"},
			wantErr:     false,
		},
		{
			name:        "postgres with database",
			serviceName: services.ServicePostgres,
			database:    "testdb",
			wantCmd:     []string{services.ClientPsql, services.User_flagPostgres, services.DefaultUserPostgres, services.Database_flagPostgres, "testdb"},
			wantErr:     false,
		},
		{
			name:        "unsupported service",
			serviceName: "unsupported",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := handler.getConnectionCommand(tt.serviceName, tt.database, tt.user, tt.host, tt.port, tt.readOnly)

			if (err != nil) != tt.wantErr {
				t.Errorf("getConnectionCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(cmd) != len(tt.wantCmd) {
					t.Errorf("getConnectionCommand() = %v, want %v", cmd, tt.wantCmd)
					return
				}

				for i, arg := range cmd {
					if arg != tt.wantCmd[i] {
						t.Errorf("getConnectionCommand()[%d] = %v, want %v", i, arg, tt.wantCmd[i])
					}
				}
			}
		})
	}
}

func TestNewConnectHandler(t *testing.T) {
	handler := NewConnectHandler()
	if handler == nil {
		t.Error("NewConnectHandler() returned nil")
	}
}
