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
	"regexp"
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
		ctx := context.Background()
		for i := range args {
			namespace := &v1.Namespace{}
			namespace.Name = args[i]
			prompt := fmt.Sprintf("clean namespace %s in this cluster, continue clean (y/n)?", args[i])
			if confirm(prompt) {
				deleteNS, err := client.CoreV1().Namespaces().Get(ctx, args[i], v12.GetOptions{})
				if err == nil {
					if deleteNS.ObjectMeta.DeletionTimestamp.IsZero() {
						e := client.CoreV1().Namespaces().Delete(context.Background(), args[i], v12.DeleteOptions{})
						if e != nil {
							logger.Error("delete namespace error: %v", e)
						} else {
							logger.Info("delete namespace success %s", namespace.Name)
						}
					} else {
						_, e := client.CoreV1().Namespaces().Finalize(context.Background(), namespace, v12.UpdateOptions{})
						if e != nil {
							logger.Error("finalize namespace error: %v", e)
						} else {
							logger.Info("wait namespace gc %s", namespace.Name)
						}
					}
				} else {
					logger.Warn("get namespace %s is error [%s], skip clean this namespace", args[i], err.Error())
				}
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

var yesRx = regexp.MustCompile("^(?i:y(?:es)?)$")

// like y|yes|Y|YES return true
func getConfirmResult(str string) bool {
	return yesRx.MatchString(str)
}

// send the prompt and get result
func confirm(prompt string) bool {
	var (
		inputStr string
		err      error
	)
	_, err = fmt.Fprint(os.Stdout, prompt)
	if err != nil {
		logger.Error("fmt.Fprint err", err)
		os.Exit(-1)
	}

	_, err = fmt.Scanf("%s", &inputStr)
	if err != nil {
		logger.Error("fmt.Scanf err", err)
		os.Exit(-1)
	}

	return getConfirmResult(inputStr)
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
