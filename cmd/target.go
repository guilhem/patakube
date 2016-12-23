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
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// targetCmd represents the target command
var targetCmd = &cobra.Command{
	Use:   "target",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: target,
}
var c = cache.New(5*time.Minute, 30*time.Second)

func init() {
	RootCmd.AddCommand(targetCmd)
}

func target(cmd *cobra.Command, args []string) {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(c.Items()); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	})

	r.Post("/", postRoot)

	if err := http.ListenAndServe(":"+strconv.Itoa(viper.GetInt("port")), r); err != nil {
		log.Fatal(err)
	}
}

func postRoot(w http.ResponseWriter, r *http.Request) {
	var p player

	accuracy := 10
	velocity := 5

	if r.Body == nil {
		http.Error(w, "Please send a request body", 412)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if r.Header.Get("patator") != "" {
		accuracy = 5
		w.Header().Set("X-Patator", "YES!")
	} else {
		w.Header().Set("X-Patator", "no, try `player patator` command")
	}

	_, found := c.Get(p.ID)
	if found {
		rand.Seed(time.Now().Unix())

		processTime := time.Duration(velocity) * time.Second

		select {
		case <-r.Context().Done():
			return

		case <-time.After(processTime):
			// The above channel simulates some hard work.
		}

		if rand.Intn(accuracy) != 0 {
			http.Error(w, "Missed", 410)
			return
		}

		if err := c.Increment(p.ID, 1); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	} else {
		c.Set(p.ID, 1, cache.NoExpiration)
	}

	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write([]byte("Hit")); err != nil {
		log.Fatal(err)
	}
}
