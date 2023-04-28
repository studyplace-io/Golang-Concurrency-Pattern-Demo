package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/spf13/pflag"
	"golanglearning/new_project/Golang-Concurrency-Pattern-Demo/exit-gracefully-mode/signals"
	"k8s.io/klog/v2"
	"log"
	"net/http"
	"os"
	"strings"
)

// Options 模拟配置文件
type Options struct {
	Port  int
	test  string
	test1 bool
	test2 int
}

func NewOptions() *Options {
	return &Options{}
}

const (
	DefaultTest = "test"
	DefaultPort = 8080
)

// AddFlags 加入命令行参数
func (o *Options) AddFlags(flags *pflag.FlagSet) {
	flags.IntVar(&o.Port, "port", DefaultPort, "xxx")
	flags.StringVar(&o.test, "test", DefaultTest, "xxx")
	flags.IntVar(&o.test2, "test2", 1000, "xxx")
	flags.BoolVar(&o.test1, "test1", false, "xxx")

	o.addKlogFlags(flags)

}

func (o *Options) addKlogFlags(flags *pflag.FlagSet) {
	klogFlags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	klog.InitFlags(klogFlags)

	klogFlags.VisitAll(func(f *flag.Flag) {
		f.Name = fmt.Sprintf("klog-%s", strings.ReplaceAll(f.Name, "_", "-"))
	})
	flags.AddGoFlagSet(klogFlags)
}

func test(w http.ResponseWriter, req *http.Request) {
	fmt.Println("test server")
}

// Run 执行
func Run(ctx context.Context, opts *Options) {
	http.HandleFunc("/test", test)
	http.ListenAndServe(fmt.Sprintf(":%v", opts.Port), nil)
}

func Example1() {

	opts := NewOptions()
	flags := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	// 把命令行参数转为options
	opts.AddFlags(flags)
	flags.Parse(os.Args[1:])
	flags.VisitAll(func(f *pflag.Flag) {
		log.Printf("Flag: %v=%v\n", f.Name, f.Value.String())
	})

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-signals.SetupSignalHandler()
		cancel()
	}()

	// 执行server逻辑
	Run(ctx, opts)

}
