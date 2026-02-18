package testhelpers

import (
	"context"
	"errors"

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

	// Call tracking
	upCalls      int
	downCalls    int
	stopCalls    int
	startCalls   int
	restartCalls int
	psCalls      int
	logsCalls    int
	execCalls    int
}

// Builder methods
func NewMockCompose() *MockComposeAPI {
	return &MockComposeAPI{}
}

func (m *MockComposeAPI) WithUpSuccess() *MockComposeAPI {
	m.UpFunc = func(ctx context.Context, project *types.Project, options api.UpOptions) error {
		return nil
	}
	return m
}

func (m *MockComposeAPI) WithUpError(err error) *MockComposeAPI {
	m.UpFunc = func(ctx context.Context, project *types.Project, options api.UpOptions) error {
		return err
	}
	return m
}

func (m *MockComposeAPI) WithDownSuccess() *MockComposeAPI {
	m.DownFunc = func(ctx context.Context, projectName string, options api.DownOptions) error {
		return nil
	}
	return m
}

func (m *MockComposeAPI) WithDownError(err error) *MockComposeAPI {
	m.DownFunc = func(ctx context.Context, projectName string, options api.DownOptions) error {
		return err
	}
	return m
}

func (m *MockComposeAPI) WithStopSuccess() *MockComposeAPI {
	m.StopFunc = func(ctx context.Context, projectName string, options api.StopOptions) error {
		return nil
	}
	return m
}

func (m *MockComposeAPI) WithStopError(err error) *MockComposeAPI {
	m.StopFunc = func(ctx context.Context, projectName string, options api.StopOptions) error {
		return err
	}
	return m
}

func (m *MockComposeAPI) WithPsResult(containers []api.ContainerSummary) *MockComposeAPI {
	m.PsFunc = func(ctx context.Context, projectName string, options api.PsOptions) ([]api.ContainerSummary, error) {
		return containers, nil
	}
	return m
}

func (m *MockComposeAPI) WithPsError(err error) *MockComposeAPI {
	m.PsFunc = func(ctx context.Context, projectName string, options api.PsOptions) ([]api.ContainerSummary, error) {
		return nil, err
	}
	return m
}

func (m *MockComposeAPI) WithLogsSuccess() *MockComposeAPI {
	m.LogsFunc = func(ctx context.Context, projectName string, consumer api.LogConsumer, options api.LogOptions) error {
		return nil
	}
	return m
}

func (m *MockComposeAPI) WithExecSuccess() *MockComposeAPI {
	m.ExecFunc = func(ctx context.Context, projectName string, options api.RunOptions) (int, error) {
		return 0, nil
	}
	return m
}

// Call tracking methods
func (m *MockComposeAPI) UpCallCount() int      { return m.upCalls }
func (m *MockComposeAPI) DownCallCount() int    { return m.downCalls }
func (m *MockComposeAPI) StopCallCount() int    { return m.stopCalls }
func (m *MockComposeAPI) StartCallCount() int   { return m.startCalls }
func (m *MockComposeAPI) RestartCallCount() int { return m.restartCalls }
func (m *MockComposeAPI) PsCallCount() int      { return m.psCalls }
func (m *MockComposeAPI) LogsCallCount() int    { return m.logsCalls }
func (m *MockComposeAPI) ExecCallCount() int    { return m.execCalls }

func (m *MockComposeAPI) Up(ctx context.Context, project *types.Project, options api.UpOptions) error {
	m.upCalls++
	if m.UpFunc != nil {
		return m.UpFunc(ctx, project, options)
	}
	return nil
}

func (m *MockComposeAPI) Down(ctx context.Context, projectName string, options api.DownOptions) error {
	m.downCalls++
	if m.DownFunc != nil {
		return m.DownFunc(ctx, projectName, options)
	}
	return nil
}

func (m *MockComposeAPI) Stop(ctx context.Context, projectName string, options api.StopOptions) error {
	m.stopCalls++
	if m.StopFunc != nil {
		return m.StopFunc(ctx, projectName, options)
	}
	return nil
}

func (m *MockComposeAPI) Start(ctx context.Context, projectName string, options api.StartOptions) error {
	m.startCalls++
	if m.StartFunc != nil {
		return m.StartFunc(ctx, projectName, options)
	}
	return nil
}

func (m *MockComposeAPI) Restart(ctx context.Context, projectName string, options api.RestartOptions) error {
	m.restartCalls++
	if m.RestartFunc != nil {
		return m.RestartFunc(ctx, projectName, options)
	}
	return nil
}

func (m *MockComposeAPI) Ps(ctx context.Context, projectName string, options api.PsOptions) ([]api.ContainerSummary, error) {
	m.psCalls++
	if m.PsFunc != nil {
		return m.PsFunc(ctx, projectName, options)
	}
	return []api.ContainerSummary{}, nil
}

func (m *MockComposeAPI) Logs(ctx context.Context, projectName string, consumer api.LogConsumer, options api.LogOptions) error {
	m.logsCalls++
	if m.LogsFunc != nil {
		return m.LogsFunc(ctx, projectName, consumer, options)
	}
	return nil
}

func (m *MockComposeAPI) Exec(ctx context.Context, projectName string, options api.RunOptions) (int, error) {
	m.execCalls++
	if m.ExecFunc != nil {
		return m.ExecFunc(ctx, projectName, options)
	}
	return 0, nil
}

// Preset scenarios
func NewMockComposeSuccess() *MockComposeAPI {
	return NewMockCompose().
		WithUpSuccess().
		WithDownSuccess().
		WithStopSuccess()
}

func NewMockComposeWithUpError() *MockComposeAPI {
	return NewMockCompose().WithUpError(errors.New("up failed"))
}

func NewMockComposeWithDownError() *MockComposeAPI {
	return NewMockCompose().WithDownError(errors.New("down failed"))
}

func NewMockComposeWithRunningContainers() *MockComposeAPI {
	return NewMockCompose().WithPsResult([]api.ContainerSummary{
		{Name: "postgres", State: "running"},
		{Name: "redis", State: "running"},
	})
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

func (m *MockComposeAPI) Publish(ctx context.Context, project *types.Project, repository string, options api.PublishOptions) error {
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
	return nil, nil
}

func (m *MockComposeAPI) Events(ctx context.Context, projectName string, options api.EventsOptions) error {
	return nil
}

func (m *MockComposeAPI) Port(ctx context.Context, projectName string, service string, port uint16, options api.PortOptions) (string, int, error) {
	return "", 0, nil
}

func (m *MockComposeAPI) Images(ctx context.Context, projectName string, options api.ImagesOptions) (map[string]api.ImageSummary, error) {
	return nil, nil
}

func (m *MockComposeAPI) Watch(ctx context.Context, project *types.Project, options api.WatchOptions) error {
	return nil
}

func (m *MockComposeAPI) Kill(ctx context.Context, projectName string, options api.KillOptions) error {
	return nil
}

func (m *MockComposeAPI) Remove(ctx context.Context, projectName string, options api.RemoveOptions) error {
	return nil
}

func (m *MockComposeAPI) RunOneOffContainer(ctx context.Context, project *types.Project, opts api.RunOptions) (int, error) {
	return 0, nil
}

func (m *MockComposeAPI) Scale(ctx context.Context, project *types.Project, options api.ScaleOptions) error {
	return nil
}

func (m *MockComposeAPI) Attach(ctx context.Context, projectName string, options api.AttachOptions) error {
	return nil
}

func (m *MockComposeAPI) Config(ctx context.Context, project *types.Project, options api.ConfigOptions) ([]byte, error) {
	return nil, nil
}

func (m *MockComposeAPI) Generate(ctx context.Context, options api.GenerateOptions) (*types.Project, error) {
	return nil, nil
}

func (m *MockComposeAPI) List(ctx context.Context, options api.ListOptions) ([]api.Stack, error) {
	return nil, nil
}

func (m *MockComposeAPI) LoadProject(ctx context.Context, options api.ProjectLoadOptions) (*types.Project, error) {
	return nil, nil
}

func (m *MockComposeAPI) Volumes(ctx context.Context, projectName string, options api.VolumesOptions) ([]api.VolumesSummary, error) {
	return nil, nil
}

func (m *MockComposeAPI) Commit(ctx context.Context, projectName string, options api.CommitOptions) error {
	return nil
}

func (m *MockComposeAPI) Export(ctx context.Context, projectName string, options api.ExportOptions) error {
	return nil
}

func (m *MockComposeAPI) Viz(ctx context.Context, project *types.Project, options api.VizOptions) (string, error) {
	return "", nil
}

func (m *MockComposeAPI) Wait(ctx context.Context, projectName string, options api.WaitOptions) (int64, error) {
	return 0, nil
}

// MockProjectLoader is a mock implementation for project loading
type MockProjectLoader struct {
	LoadFunc func(projectName string) (*types.Project, error)
}

func NewMockProjectLoader() *MockProjectLoader {
	return &MockProjectLoader{}
}

func (m *MockProjectLoader) WithLoadSuccess(projectName string) *MockProjectLoader {
	m.LoadFunc = func(name string) (*types.Project, error) {
		return &types.Project{Name: projectName}, nil
	}
	return m
}

func (m *MockProjectLoader) WithLoadError(err error) *MockProjectLoader {
	m.LoadFunc = func(name string) (*types.Project, error) {
		return nil, err
	}
	return m
}

func (m *MockProjectLoader) Load(projectName string) (*types.Project, error) {
	if m.LoadFunc != nil {
		return m.LoadFunc(projectName)
	}
	return &types.Project{Name: projectName}, nil
}

// MockLogConsumer is a mock implementation for log consumption
type MockLogConsumer struct {
	Logs []string
}

func (m *MockLogConsumer) Log(service, container, message string) {
	m.Logs = append(m.Logs, message)
}

func (m *MockLogConsumer) Status(container, status string) {}

func (m *MockLogConsumer) Register(container string) {}
