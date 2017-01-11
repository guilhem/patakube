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
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"encoding/json"

	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// playerCmd represents the player command
var playerCmd = &cobra.Command{
	Use:   "player",
	Short: "Launch a patakube player",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: playerRun,
}

func init() {
	RootCmd.AddCommand(playerCmd)

	playerCmd.PersistentFlags().String("target", "target.patakube.svc.cluster.local", "target to shoot")
	if err := viper.BindPFlag("target", playerCmd.PersistentFlags().Lookup("target")); err != nil {
		log.Fatal(err)
	}

	playerCmd.Flags().Bool("patator", false, "Use a patator")
	if err := viper.BindPFlag("patator", playerCmd.Flags().Lookup("patator")); err != nil {
		log.Fatal(err)
	}

	playerCmd.Flags().Int("patator-port", 8081, "Port to reach patator")
	if err := viper.BindPFlag("patator-port", playerCmd.Flags().Lookup("patator-port")); err != nil {
		log.Fatal(err)
	}
}

type player struct {
	ID         string
	Namespace  string
	ClusterURL string
}

func playerRun(cmd *cobra.Command, args []string) {
	p := player{ID: os.Getenv("NAME")}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(fmt.Sprintf("Player %s \n", p.ID))); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	})

	// Slow handlers/operations.
	r.Group(func(r chi.Router) {
		// Stop processing after 30 seconds.
		r.Use(middleware.Timeout(30 * time.Second))

		// Only one request will be processed at a time.
		r.Use(middleware.Throttle(1))

		r.Post("/potato", func(w http.ResponseWriter, r *http.Request) {
			b := new(bytes.Buffer)
			if err := json.NewEncoder(b).Encode(p); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			target := "http://" + viper.GetString("target")
			if viper.GetBool("patator") {
				target = "http://localhost:" + strconv.Itoa(viper.GetInt("patator-port"))
			}
			res, err := http.Post(target, "application/json", b)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			if res.StatusCode != http.StatusCreated {
				http.Error(w, "Missed", 410)
				return
			}

			w.WriteHeader(http.StatusCreated)
			if _, err := w.Write([]byte("Hit")); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
		})
	})

	if err := http.ListenAndServe(":"+strconv.Itoa(viper.GetInt("port")), r); err != nil {
		log.Fatal(err)
	}
}
