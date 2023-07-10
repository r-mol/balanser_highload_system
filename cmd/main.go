package main

import (
	"fmt"
	"github.com/r-mol/balanser_highload_system/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/http"
	"os"
)

func Main(configPath, address string) error {
	cnf, err := config.ParseBalancerConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to parse Balancer config: %w", err)
	}

	b, err := config.GetBalancerFromConfig(cnf)
	if err != nil {
		return fmt.Errorf("failed to get balancer from config: %w", err)
	}

	log.Infoln("balancer started at address: " + address)
	return http.ListenAndServe(address, b)
}

func GetStarterCmd() *cobra.Command {
	var configPath, address string

	cmd := &cobra.Command{
		Use:     "start",
		Version: "0.0.1",
		Short:   "launch load balancer",
		Run: func(cmd *cobra.Command, args []string) {
			err := Main(configPath, address)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "", "path to config file")
	cmd.Flags().StringVar(&address, "address", "", "address of balancer")
	_ = cmd.MarkFlagRequired("config")
	_ = cmd.MarkFlagRequired("address")

	return cmd
}

func main() {
	if err := GetStarterCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
