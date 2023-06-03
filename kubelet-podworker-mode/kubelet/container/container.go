package container

import "k8s.io/klog/v2"

// CRI 容器运行时接口
type CRI interface {
	// RunPodSandbox 运行沙箱容器
	RunPodSandbox() error
	// StopPodSandbox 停止沙箱容器
	StopPodSandbox() error
	// RemovePodSandbox 删除沙箱容器
	RemovePodSandbox() error
	// CreateContainer 创建容器
	CreateContainer() error
	// StartContainer 启动容器
	StartContainer() error
	// StopContainer 停止容器
	StopContainer() error
	// RemoveContainer 删除容器
	RemoveContainer() error
	// PullImage 拉取镜像
	PullImage(string) error
	// RemoveImage 删除镜像
	RemoveImage(string) error
}

var _ CRI = &Container{}

// Container 容器对象
type Container struct {
	Name            string
	Image           string
	ImagePullPolicy string
	Env             map[string]string
	Status          Status
}

type Status string

const (
	NoRunning       Status = "noRunning"
	Running         Status = "running"
	Fail            Status = "fail"
	Creating        Status = "creating"
	Removing        Status = "removing"
	Stopping        Status = "stopping"
	ImagePulling    Status = "imagePulling"
	ImageRemoving   Status = "imageRemoving"
	SandboxRunning  Status = "sandboxRunning"
	SandboxStopping Status = "sandboxStopping"
	SandboxRemoving Status = "sandboxRemoving"
)

func (c *Container) RunPodSandbox() error {
	klog.Infof("run pod sandbox %s", c.Name)
	c.Status = SandboxRunning
	return nil
}

func (c *Container) StopPodSandbox() error {
	klog.Infof("stop pod sandbox %s", c.Name)
	c.Status = SandboxStopping
	return nil
}

func (c *Container) RemovePodSandbox() error {
	klog.Infof("remove pod sandbox %s", c.Name)
	c.Status = SandboxRemoving
	return nil
}

func (c *Container) CreateContainer() error {
	klog.Infof("create container %s", c.Name)
	c.Status = Creating
	return nil
}

func (c *Container) StartContainer() error {
	klog.Infof("start container %s", c.Name)
	c.Status = Running
	return nil
}

func (c *Container) StopContainer() error {
	klog.Infof("stop container %s", c.Name)
	c.Status = Stopping
	return nil
}

func (c *Container) RemoveContainer() error {
	klog.Infof("remove container %s", c.Name)
	c.Status = Removing
	return nil
}

func (c *Container) PullImage(s string) error {
	klog.Infof("pull image %s for container %s", s, c.Name)
	c.Status = ImagePulling
	return nil
}

func (c *Container) RemoveImage(s string) error {
	klog.Infof("remove image %s for container %s", s, c.Name)
	c.Status = ImageRemoving
	return nil
}