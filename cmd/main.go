package main

import (
	"fmt"
	"github.com/r-mol/balanser_highload_system/config"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

func Main(configPath string) error {
	cnf, err := config.ParseBalancerConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to parse Balancer config: %w", err)
	}

	b, err := config.GetBalancerFromConfig(cnf)
	if err != nil {
		return fmt.Errorf("failed to get balancer from config: %w", err)
	}

	fmt.Println("balancer started at port 127.0.0.1:8080")
	return http.ListenAndServe("127.0.0.1:8080", b)
}

func GetStarterCmd() *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:     "start",
		Version: "0.0.1",
		Short:   "launch load balancer",
		Run: func(cmd *cobra.Command, args []string) {
			err := Main(configPath)
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "", "path to config file")
	_ = cmd.MarkFlagRequired("config")

	return cmd
}

func main() {
	if err := GetStarterCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
