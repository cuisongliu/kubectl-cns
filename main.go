/*
Copyright © 2021 cuisongliu@qq.com

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
package main

import (
	"context"
	"fmt"
	"github.com/cuisongliu/clean-ns/logger"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	hd "k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "clean-ns",
	Short: "clean namespace",
	Run: func(cmd *cobra.Command, args []string) {
		config := getRestConfig()
		if config == nil {
			logger.Fatal("kubernetes config is not config")
		}
		client := kubernetes.NewForConfigOrDie(config)

		for i := range args {
			namespace := &v1.Namespace{}
			namespace.Name = args[i]
			_, err := client.CoreV1().Namespaces().Finalize(context.Background(), namespace, v12.UpdateOptions{})
			if err != nil {
				logger.Error("finalize namespace error: %v", err)
			} else {
				logger.Error("wait namespace gc %s", namespace.Name)
			}
		}
	},
}

func getRestConfig() *rest.Config {
	config, err := rest.InClusterConfig()
	if err != nil {
		var kubeconfig = filepath.Join(hd.HomeDir(), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil
		}
	}
	return config
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
