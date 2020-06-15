/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/potterhe/nodelocaldns-injector/pkg/inject"
	"github.com/spf13/cobra"
)

var (
	certFile string
	keyFile  string
	port     int
)

// webhookCmd represents the webhook command
var webhookCmd = &cobra.Command{
	Use:   "webhook",
	Short: "Starts a HTTP server, useful for inject dnsConfig MutatingAdmissionWebhook",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("webhook called")

		parameters := inject.WebhookParameters{
			CertFile: certFile,
			KeyFile:  keyFile,
			Port:     port,
		}

		wh, err := inject.NewWebhook(parameters)
		if err != nil {
			panic(err)
		}

		err = wh.Serve()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(webhookCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// webhookCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// webhookCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	webhookCmd.Flags().StringVar(&certFile, "tls-cert-file", "",
		"File containing the default x509 Certificate for HTTPS. (CA cert, if any, concatenated after server cert).")
	webhookCmd.Flags().StringVar(&keyFile, "tls-private-key-file", "",
		"File containing the default x509 private key matching --tls-cert-file.")
	webhookCmd.Flags().IntVar(&port, "port", 443,
		"Secure port that the webhook listens on")
}
