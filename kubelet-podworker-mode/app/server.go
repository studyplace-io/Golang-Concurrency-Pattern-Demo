package app

import "github.com/spf13/cobra"

func NewKubeletCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "kubelet",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: 未完成
			return nil
		},
	}
	return cmd
}
