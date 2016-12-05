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
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// patatorCmd represents the patator command
var patatorCmd = &cobra.Command{
	Use:   "patator",
	Short: "A potato launcher",
	Long: `This tool help you to improve aiming of your shoots.

	Be careful, this tool can be quite capricious ;)`,
	Run: patator,
}

func init() {
	playerCmd.AddCommand(patatorCmd)
}

func patator(cmd *cobra.Command, args []string) {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			http.Error(w, "Please send a request body", 412)
			return
		}

		w.Header().Set("patator", "true")

		req, err := http.NewRequest("POST", "http://"+viper.GetString("target"), r.Body)
		req.Header.Set("patator", "true")
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()

		responseBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		w.WriteHeader(res.StatusCode)
		w.Write(responseBody)
	})

	http.ListenAndServe(":"+strconv.Itoa(viper.GetInt("port")), r)
}
