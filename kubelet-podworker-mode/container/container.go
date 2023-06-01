package container

type CRI interface {

	// 沙箱
	RunPodSandbox() error
	StopPodSandbox() error
	RemovePodSandbox() error

	// 容器
	CreateContainer() error
	StartContainer() error
	StopContainer() error
	RemoveContainer() error

	// 镜像
	PullImage(string) error
	RemoveImage(string) error
}

var _ CRI = &Container{}

type Container struct {
	Name            string
	Image           string
	ImagePullPolicy string
	Env             map[string]string
	Status          Status
}

type Status string

const (
	Running Status = "running"
	Fail    Status = "fail"
)

func (c *Container) RunPodSandbox() error {
	panic("implement me")
}

func (c *Container) StopPodSandbox() error {
	panic("implement me")
}

func (c *Container) RemovePodSandbox() error {
	panic("implement me")
}

func (c *Container) CreateContainer() error {
	panic("implement me")
}

func (c *Container) StartContainer() error {
	panic("implement me")
}

func (c *Container) StopContainer() error {
	panic("implement me")
}

func (c *Container) RemoveContainer() error {
	panic("implement me")
}

func (c *Container) PullImage(s string) error {
	panic("implement me")
}

func (c *Container) RemoveImage(s string) error {
	panic("implement me")
}

