package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/r-mol/balanser_highload_system/internal/balancer"
	data_transfer_api "github.com/r-mol/balanser_highload_system/protos"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"

	"github.com/r-mol/balanser_highload_system/config"
	"github.com/spf13/cobra"
)

func Main(configPath, address, promAddr string) error {
	cnf, err := config.ParseBalancerConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to parse Balancer config: %w", err)
	}

	// logger
	logger := log.New()
	logger.Level = log.DebugLevel

	// metrics
	metrics := balancer.NewMetrics()
	if err := metrics.Register(); err != nil {
		logger.Fatal("failed to create metrics", err)
	}
	http.Handle("/metrics", promhttp.Handler())

	go func() {
		logger.Infoln("starting metric server: ", promAddr)
		if err := http.ListenAndServe(promAddr, nil); err != nil {
			logger.Fatal(err)
		}
	}()

	// balancer
	lb, err := config.GetBalancerFromConfig(cnf, logger, metrics)
	if err != nil {
		return fmt.Errorf("failed to get balancer from config: %w", err)
	}

	lb.Logger.Infoln("balancer started at address: ", address)
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
	var configPath, address, promAddress string

	cmd := &cobra.Command{
		Use:     "start",
		Version: "0.0.1",
		Short:   "launch load balancer",
		Run: func(cmd *cobra.Command, args []string) {
			err := Main(configPath, address, promAddress)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "", "path to config file")
	cmd.Flags().StringVar(&address, "address", "", "address of balancer")
	cmd.Flags().StringVar(&promAddress, "prometheus_address", "", "address of metric server")

	_ = cmd.MarkFlagRequired("config")
	_ = cmd.MarkFlagRequired("address")
	_ = cmd.MarkFlagRequired("prometheus_address")

	return cmd
}

func main() {
	if err := GetStarterCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
