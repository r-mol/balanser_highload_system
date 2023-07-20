package main

import (
	"fmt"
	data_transfer_api "github.com/r-mol/balanser_highload_system/protos"
	"google.golang.org/grpc"
	"net"
	"os"

	"github.com/r-mol/balanser_highload_system/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Main(configPath, address string) error {
	cnf, err := config.ParseBalancerConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to parse Balancer config: %w", err)
	}

	lb, err := config.GetBalancerFromConfig(cnf)
	if err != nil {
		return fmt.Errorf("failed to get balancer from config: %w", err)
	}

	log.Infoln("balancer started at address: " + address)
	lis, err := net.Listen("tcp4", address)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	defer lis.Close()

	grpcServer := grpc.NewServer()
	data_transfer_api.RegisterKeyValueServiceServer(grpcServer, lb)

	err = grpcServer.Serve(lis)
	if err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

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
