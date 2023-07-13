package main

import (
	"fmt"
	"github.com/r-mol/balanser_highload_system/internal/handler"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/r-mol/balanser_highload_system/config"
	"github.com/spf13/cobra"
)

func Main(configPath, address string) error {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	cnf, err := config.ParseBalancerConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to parse LoadBalancer config: %w", err)
	}

	lb, err := config.GetBalancerFromConfig(cnf)
	if err != nil {
		return fmt.Errorf("failed to get balancer from config: %w", err)
	}

	e.Use(handler.MiddlewareBalancer(lb))

	e.Logger.Fatal(e.Start(address))

	return nil
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
