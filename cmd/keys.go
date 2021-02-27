// Copyright © 2018 Jack Zampolin <jack@blockstack.com>
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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/twjang/keyserver/api"
)

// versionCmd represents the version command
var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Runs keys calls",
}

// /keys/list POST
var keysList = &cobra.Command{
	Use:   "list [password]",
	Args:  cobra.ExactArgs(1),
	Short: "Fetch all keys managed by the keyserver",
	Run: func(cmd *cobra.Command, args []string) {
		url := fmt.Sprintf("http://localhost:%d/keys/list", server.Port)
		req := api.GetKeyBody{
			Password: args[0],
		}
		reqbytes, err := json.Marshal(req)
		if err != nil {
			panic(err)
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqbytes))
		if err != nil {
			log.Fatalf("error fetching %s", url)
			return
		}
		if resp.StatusCode != 200 {
			log.Fatalf("non 200 respose code")
			return
		}
		out, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("failed reading response body")
			return
		}
		fmt.Println(string(out))
	},
}

// /keys/create POST
var keysPost = &cobra.Command{
	Use:   "create [password] [name] [mnemonic]",
	Args:  cobra.RangeArgs(2, 3),
	Short: "Add a new key to the keyserver, optionally pass a mnemonic to restore the key",
	Run: func(cmd *cobra.Command, args []string) {
		url := fmt.Sprintf("http://localhost:%d/keys/create", server.Port)
		var addNP api.AddNewKey
		if len(args) == 2 {
			addNP = api.AddNewKey{Name: args[1], Password: args[0]}
		} else if len(args) == 3 {
			addNP = api.AddNewKey{Name: args[1], Password: args[0], Mnemonic: args[2]}
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(addNP.Marshal()))
		if err != nil {
			log.Fatalf("error fetching %s", url)
			return
		}
		if resp.StatusCode != 200 {
			log.Fatalf("non 200 respose code")
			return
		}
		out, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("failed reading response body")
			return
		}
		fmt.Println(string(out))
	},
}

// /keys/get/{name} POST
var keyGet = &cobra.Command{
	Use:   "show [password] [name]",
	Args:  cobra.ExactArgs(2),
	Short: "Fetch details for one key",
	Run: func(cmd *cobra.Command, args []string) {
		req := api.GetKeyBody{
			Password: args[0],
		}
		reqbytes, err := json.Marshal(req)
		if err != nil {
			panic(err)
		}

		url := fmt.Sprintf("http://localhost:%d/keys/get/%s", server.Port, args[1])
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqbytes))

		if err != nil {
			log.Fatalf("error fetching %s", url)
			return
		}
		if resp.StatusCode != 200 {
			log.Fatalf("non 200 respose code")
			return
		}
		out, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("failed reading response body")
			return
		}
		fmt.Println(string(out))
	},
}

// /keys/delete/{name} DELETE
var keyDelete = &cobra.Command{
	Use:   "delete [password] [name]",
	Args:  cobra.ExactArgs(2),
	Short: "Delete a key",
	Run: func(cmd *cobra.Command, args []string) {
		url := fmt.Sprintf("http://localhost:%d/keys/delete/%s", server.Port, args[1])
		kb := api.DeleteKeyBody{Password: args[0]}
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(kb.Marshal()))
		if err != nil {
			log.Fatalf("error fetching %s", url)
			return
		}
		if resp.StatusCode != 200 {
			log.Fatalf("non 200 respose code %d", resp.StatusCode)
			return
		}
		out, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("failed reading response body")
			return
		}
		fmt.Println(out)
	},
}

func init() {
	keysCmd.AddCommand(keysList)
	keysCmd.AddCommand(keysPost)
	keysCmd.AddCommand(keyGet)
	keysCmd.AddCommand(keyDelete)
	rootCmd.AddCommand(keysCmd)
}
