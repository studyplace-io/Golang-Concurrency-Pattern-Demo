package app

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/study-io/Golang-Concurrency-Pattern-Demo/kube-controller-manager-mode/app/config"
	"github.com/study-io/Golang-Concurrency-Pattern-Demo/kube-controller-manager-mode/app/options"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/component-base/version/verflag"
	"k8s.io/klog/v2"
	"os"
)

// NewControllerManagerCommand creates a *cobra.Command object with default parameters
func NewControllerManagerCommand() *cobra.Command {
	// 配置文件
	s, err := options.NewKubeControllerManagerOptions()
	if err != nil {
		klog.Fatalf("unable to initialize command options: %v", err)
	}

	cmd := &cobra.Command{
		Use: "kube-controller-manager-mode",
		Long: `The Kubernetes controller manager is a daemon that embeds
the core control loops shipped with Kubernetes. In applications of robotics and
automation, a control loop is a non-terminating loop that regulates the state of
the system. In Kubernetes, a controller is a control loop that watches the shared
state of the cluster through the apiserver and makes changes attempting to move the
current state towards the desired state. Examples of controllers that ship with
Kubernetes today are the replication controller, endpoints controller, namespace
controller, and serviceaccounts controller.`,

		RunE: func(cmd *cobra.Command, args []string) error {
			verflag.PrintAndExitIfRequested()

			// 执行入口
			c, _ := s.Config()
			return Run(c.Complete(), wait.NeverStop)
		},
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		},
	}

	flags := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	s.AddFlags(flags)
	flags.Parse(os.Args[1:])
	flags.VisitAll(func(f *pflag.Flag) {
		klog.Infof("Flag: %v=%v\n", f.Name, f.Value.String())
	})

	fs := cmd.Flags()
	fs.AddFlagSet(flags)

	return cmd
}

type Interface interface {
	// Name returns the canonical name of the controller.
	Name() string
}

type InitFunc func(ctx context.Context) (controller Interface, enabled bool, err error)

type ControllerInitializersFunc func() (initializers map[string]InitFunc)

// Run runs the KubeControllerManagerOptions.
func Run(c *config.CompletedConfig, stopCh <-chan struct{}) error {

	// run 执行方法，下面会调用此func
	run := func(ctx context.Context, initializersFunc ControllerInitializersFunc) {

		// 调用各资源控制器的initializersFunc，返回的是map[string]InitFunc
		controllerInitializers := initializersFunc()
		// 启动每个controller
		if err := StartControllers(ctx, controllerInitializers); err != nil {
			klog.Fatalf("error starting controllers: %v", err)
		}

		<-ctx.Done()
	}
	ctx, _ := wait.ContextForChannel(stopCh)
	run(ctx, NewControllerInitializers)

	<-stopCh
	return nil
}

func StartControllers(ctx context.Context, controllers map[string]InitFunc) error {
	// Always start the SA token controller first using a full-power client, since it needs to mint tokens for the rest
	// If this fails, just return here and fail since other controllers won'event-mode be able to get credentials.

	// 遍历传进来的initFn，并执行
	for controllerName, initFn := range controllers {

		klog.V(1).Infof("Starting %q", controllerName)
		// 执行initFn方法，其中 controllerCtx是重要的struct，
		_, started, err := initFn(ctx)
		if err != nil {
			klog.Errorf("Error starting %q", controllerName)
			return err
		}
		if !started {
			klog.Warningf("Skipping %q", controllerName)
			continue
		}
	}

	return nil
}

func NewControllerInitializers() map[string]InitFunc {
	controllers := map[string]InitFunc{}

	// All of the controllers must have unique names, or else we will explode.
	register := func(name string, fn InitFunc) {
		if _, found := controllers[name]; found {
			panic(fmt.Sprintf("controller name %q was registered twice", name))
		}
		controllers[name] = fn
	}
	register("myController1", startMyController1)
	register("myContoller2", startMyController2)

	return controllers
}
