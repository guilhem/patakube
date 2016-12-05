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
	Run: player,
}

func init() {
	RootCmd.AddCommand(playerCmd)

	playerCmd.PersistentFlags().String("target", "target.patakube.svc.cluster.local", "target to shoot")
	viper.BindPFlag("target", playerCmd.PersistentFlags().Lookup("target"))

	playerCmd.Flags().Bool("patator", false, "Use a patator")
	viper.BindPFlag("patator", playerCmd.Flags().Lookup("patator"))

	playerCmd.Flags().Int("patator-port", 8081, "Port to reach patator")
	viper.BindPFlag("patator-port", playerCmd.Flags().Lookup("patator-port"))
}

type Player struct {
	ID string
}

func player(cmd *cobra.Command, args []string) {
	player := Player{ID: os.Getenv("NAME")}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("Player %s \n", player.ID)))
	})

	// Slow handlers/operations.
	r.Group(func(r chi.Router) {
		// Stop processing after 30 seconds.
		r.Use(middleware.Timeout(30 * time.Second))

		// Only one request will be processed at a time.
		r.Use(middleware.Throttle(1))

		r.Post("/potato", func(w http.ResponseWriter, r *http.Request) {
			b := new(bytes.Buffer)
			json.NewEncoder(b).Encode(player)
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
			w.Write([]byte("Hit"))
		})
	})

	http.ListenAndServe(":"+strconv.Itoa(viper.GetInt("port")), r)
}
