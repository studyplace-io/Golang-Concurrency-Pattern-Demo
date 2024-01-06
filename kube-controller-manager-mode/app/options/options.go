package options

import (
	"flag"
	"fmt"
	"github.com/spf13/pflag"
	"github.com/study-io/Golang-Concurrency-Pattern-Demo/kube-controller-manager-mode/app/config"
	"k8s.io/klog/v2"
	"os"
	"strings"
)

type KubeControllerManagerOptions struct {
}

// NewKubeControllerManagerOptions creates a new KubeControllerManagerOptions with a default config.
func NewKubeControllerManagerOptions() (*KubeControllerManagerOptions, error) {

	s := KubeControllerManagerOptions{}

	return &s, nil
}

func (s KubeControllerManagerOptions) Config() (*config.Config, error) {

	c := &config.Config{}

	return c, nil
}

// AddFlags 加入命令行参数
func (s *KubeControllerManagerOptions) AddFlags(flags *pflag.FlagSet) {
	s.addKlogFlags(flags)
}

func (s *KubeControllerManagerOptions) addKlogFlags(flags *pflag.FlagSet) {
	klogFlags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	klog.InitFlags(klogFlags)

	klogFlags.VisitAll(func(f *flag.Flag) {
		f.Name = fmt.Sprintf("klog-%s", strings.ReplaceAll(f.Name, "_", "-"))
	})
	flags.AddGoFlagSet(klogFlags)
}
