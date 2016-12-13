// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"net/http"
	"strconv"
	"text/template"

	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/rest"
)

var configuratorCmd = &cobra.Command{
	Use:   "configurator",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: configurator,
}

func init() {
	RootCmd.AddCommand(configuratorCmd)
}

func configurator(cmd *cobra.Command, args []string) {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/:playerID", func(w http.ResponseWriter, r *http.Request) {
		playerID := chi.URLParam(r, "playerID")
		ns := v1.Namespace{
			ObjectMeta: v1.ObjectMeta{
				Name: playerID,
			},
		}
		if _, err := clientset.Core().Namespaces().Create(&ns); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		t, err := template.ParseFiles("templates/config.sh")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		if err := t.Execute(w, playerID); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	})

	http.ListenAndServe(":"+strconv.Itoa(viper.GetInt("port")), r)
}