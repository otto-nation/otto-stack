package testhelpers

import (
	"context"
	"io"

	"github.com/compose-spec/compose-go/v2/types"
	"github.com/docker/compose/v5/pkg/api"
)

// MockComposeAPI is a mock implementation of api.Compose for testing
type MockComposeAPI struct {
	UpFunc      func(ctx context.Context, project *types.Project, options api.UpOptions) error
	DownFunc    func(ctx context.Context, projectName string, options api.DownOptions) error
	StopFunc    func(ctx context.Context, projectName string, options api.StopOptions) error
	StartFunc   func(ctx context.Context, projectName string, options api.StartOptions) error
	RestartFunc func(ctx context.Context, projectName string, options api.RestartOptions) error
	PsFunc      func(ctx context.Context, projectName string, options api.PsOptions) ([]api.ContainerSummary, error)
	LogsFunc    func(ctx context.Context, projectName string, consumer api.LogConsumer, options api.LogOptions) error
	ExecFunc    func(ctx context.Context, projectName string, options api.RunOptions) (int, error)
}

func (m *MockComposeAPI) Up(ctx context.Context, project *types.Project, options api.UpOptions) error {
	if m.UpFunc != nil {
		return m.UpFunc(ctx, project, options)
	}
	return nil
}

func (m *MockComposeAPI) Down(ctx context.Context, projectName string, options api.DownOptions) error {
	if m.DownFunc != nil {
		return m.DownFunc(ctx, projectName, options)
	}
	return nil
}

func (m *MockComposeAPI) Stop(ctx context.Context, projectName string, options api.StopOptions) error {
	if m.StopFunc != nil {
		return m.StopFunc(ctx, projectName, options)
	}
	return nil
}

func (m *MockComposeAPI) Start(ctx context.Context, projectName string, options api.StartOptions) error {
	if m.StartFunc != nil {
		return m.StartFunc(ctx, projectName, options)
	}
	return nil
}

func (m *MockComposeAPI) Restart(ctx context.Context, projectName string, options api.RestartOptions) error {
	if m.RestartFunc != nil {
		return m.RestartFunc(ctx, projectName, options)
	}
	return nil
}

func (m *MockComposeAPI) Ps(ctx context.Context, projectName string, options api.PsOptions) ([]api.ContainerSummary, error) {
	if m.PsFunc != nil {
		return m.PsFunc(ctx, projectName, options)
	}
	return []api.ContainerSummary{}, nil
}

func (m *MockComposeAPI) Logs(ctx context.Context, projectName string, consumer api.LogConsumer, options api.LogOptions) error {
	if m.LogsFunc != nil {
		return m.LogsFunc(ctx, projectName, consumer, options)
	}
	return nil
}

func (m *MockComposeAPI) Exec(ctx context.Context, projectName string, options api.RunOptions) (int, error) {
	if m.ExecFunc != nil {
		return m.ExecFunc(ctx, projectName, options)
	}
	return 0, nil
}

// Implement remaining required methods with no-ops
func (m *MockComposeAPI) Build(ctx context.Context, project *types.Project, options api.BuildOptions) error {
	return nil
}

func (m *MockComposeAPI) Push(ctx context.Context, project *types.Project, options api.PushOptions) error {
	return nil
}

func (m *MockComposeAPI) Pull(ctx context.Context, project *types.Project, options api.PullOptions) error {
	return nil
}

func (m *MockComposeAPI) Create(ctx context.Context, project *types.Project, options api.CreateOptions) error {
	return nil
}

func (m *MockComposeAPI) Copy(ctx context.Context, projectName string, options api.CopyOptions) error {
	return nil
}

func (m *MockComposeAPI) Pause(ctx context.Context, projectName string, options api.PauseOptions) error {
	return nil
}

func (m *MockComposeAPI) UnPause(ctx context.Context, projectName string, options api.PauseOptions) error {
	return nil
}

func (m *MockComposeAPI) Top(ctx context.Context, projectName string, services []string) ([]api.ContainerProcSummary, error) {
	return []api.ContainerProcSummary{}, nil
}

func (m *MockComposeAPI) Events(ctx context.Context, projectName string, options api.EventsOptions) error {
	return nil
}

func (m *MockComposeAPI) Port(ctx context.Context, projectName string, service string, port uint16, options api.PortOptions) (string, int, error) {
	return "", 0, nil
}

func (m *MockComposeAPI) Images(ctx context.Context, projectName string, options api.ImagesOptions) (map[string]api.ImageSummary, error) {
	return map[string]api.ImageSummary{}, nil
}

func (m *MockComposeAPI) Watch(ctx context.Context, project *types.Project, options api.WatchOptions) error {
	return nil
}

func (m *MockComposeAPI) MaxConcurrency(parallel int) {
}

func (m *MockComposeAPI) DryRunMode(ctx context.Context, dryRun bool) (context.Context, error) {
	return ctx, nil
}

func (m *MockComposeAPI) Viz(ctx context.Context, project *types.Project, options api.VizOptions) (string, error) {
	return "", nil
}

func (m *MockComposeAPI) Wait(ctx context.Context, projectName string, options api.WaitOptions) (int64, error) {
	return 0, nil
}

func (m *MockComposeAPI) Publish(ctx context.Context, project *types.Project, repository string, options api.PublishOptions) error {
	return nil
}

func (m *MockComposeAPI) Scale(ctx context.Context, project *types.Project, options api.ScaleOptions) error {
	return nil
}

func (m *MockComposeAPI) Kill(ctx context.Context, projectName string, options api.KillOptions) error {
	return nil
}

func (m *MockComposeAPI) Remove(ctx context.Context, projectName string, options api.RemoveOptions) error {
	return nil
}

func (m *MockComposeAPI) RunOneOffContainer(ctx context.Context, project *types.Project, options api.RunOptions) (int, error) {
	return 0, nil
}

func (m *MockComposeAPI) Attach(ctx context.Context, projectName string, options api.AttachOptions) error {
	return nil
}

func (m *MockComposeAPI) Config(ctx context.Context, project *types.Project, options api.ConfigOptions) ([]byte, error) {
	return []byte{}, nil
}

func (m *MockComposeAPI) Interceptors() []any {
	return []any{}
}

func (m *MockComposeAPI) WithInterceptors(interceptors ...any) api.Compose {
	return m
}

func (m *MockComposeAPI) Commit(ctx context.Context, projectName string, options api.CommitOptions) error {
	return nil
}

func (m *MockComposeAPI) Export(ctx context.Context, projectName string, options api.ExportOptions) error {
	return nil
}

func (m *MockComposeAPI) Generate(ctx context.Context, options api.GenerateOptions) (*types.Project, error) {
	return &types.Project{}, nil
}

func (m *MockComposeAPI) List(ctx context.Context, options api.ListOptions) ([]api.Stack, error) {
	return []api.Stack{}, nil
}

func (m *MockComposeAPI) LoadProject(ctx context.Context, options api.ProjectLoadOptions) (*types.Project, error) {
	return &types.Project{}, nil
}

func (m *MockComposeAPI) Volumes(ctx context.Context, projectName string, options api.VolumesOptions) ([]api.VolumesSummary, error) {
	return []api.VolumesSummary{}, nil
}

// MockProjectLoader is a mock implementation of ProjectLoader
type MockProjectLoader struct {
	LoadFunc func(projectName string) (*types.Project, error)
}

func (m *MockProjectLoader) Load(projectName string) (*types.Project, error) {
	if m.LoadFunc != nil {
		return m.LoadFunc(projectName)
	}
	return &types.Project{Name: projectName}, nil
}

// MockLogConsumer implements api.LogConsumer for testing
type MockLogConsumer struct {
	LogFunc func(service, container, message string)
}

func (m *MockLogConsumer) Log(service, container, message string) {
	if m.LogFunc != nil {
		m.LogFunc(service, container, message)
	}
}

func (m *MockLogConsumer) Status(container, msg string) {
}

func (m *MockLogConsumer) Register(name string) {
}

func (m *MockLogConsumer) Err(service, container string, err error) {
}

func (m *MockLogConsumer) Stdout(service, container string, message string) {
}

func (m *MockLogConsumer) Stderr(service, container string, message string) {
}

func (m *MockLogConsumer) GetWriter(container string) io.Writer {
	return io.Discard
}
