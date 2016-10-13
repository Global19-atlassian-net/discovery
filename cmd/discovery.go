package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/quantum/discovery/handlers"

	"github.com/coreos/etcd/client"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	DefaultTimeout = 5 * time.Second
)

var (
	etcdURLs      string //comma separated list of etcd urls
	discoveryAddr string
	host          string
	dnsAddr       string
)

// RootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "discovery",
	Short: "discovery tool to bootstrap etcd nodes and help clients find them.",
	Long:  `https://github.com/quantum/discovery`,
	Run: func(cmd *cobra.Command, args []string) {
		startDiscovery()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(
		&etcdURLs,
		"etcd-urls",
		"http://127.0.0.1:2379",
		"comma separated list of etcd listen URLs")
	rootCmd.PersistentFlags().StringVar(
		&discoveryAddr,
		"addr",
		"http://127.0.0.1:8087",
		"address:port of the discovery service")
	rootCmd.PersistentFlags().StringVar(
		&discoveryAddr,
		"host",
		"https://discovery.castle.io",
		"the host url to prepend to /new requests.")
	// load the environment variables
	setFlagsFromEnv(rootCmd.Flags())
	rootCmd.MarkFlagRequired("etcd-urls")
}

func startDiscovery() {
	etcdMachines := strings.Split(etcdURLs, ",")
	cfg := client.Config{
		Endpoints: etcdMachines,
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: DefaultTimeout,
	}
	c, err := client.New(cfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	etcdClient := client.NewKeysAPI(c)

	// FIXME: remove the next line after pointing the discovery to discovery.castle.io
	host = discoveryAddr

	handlers.Setup(host)
	h := handlers.NewHandler(etcdClient)
	r := handlers.NewRouter(h.GetRoutes())

	u, err := url.Parse(discoveryAddr)
	if err != nil {
		log.Printf("error in parsing url: %+v", err)
	}
	if err := http.ListenAndServe(u.Host, r); err != nil {
		log.Printf("API error: %+v", err)
	}
}

func setFlagsFromEnv(flags *pflag.FlagSet) error {
	flags.VisitAll(func(f *pflag.Flag) {
		envVar := "DISCOVERY_" + strings.Replace(strings.ToUpper(f.Name), "-", "_", -1)
		value := os.Getenv(envVar)
		if value != "" {
			// Set the environment variable. Will override default values, but be overridden by command line parameters.
			flags.Set(f.Name, value)
		}
	})

	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
