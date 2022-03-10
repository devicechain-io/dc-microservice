/*
Copyright © 2022 SiteWhere LLC - All Rights Reserved
Unauthorized copying of this file, via any medium is strictly prohibited.
Proprietary and confidential.
*/

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/devicechain-io/dc-microservice/cmd"
	"github.com/devicechain-io/dc-microservice/core"
	"github.com/devicechain-io/dc-microservice/generator"
	"github.com/fatih/color"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	testGenerate()
	banner()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	cmd.Execute()
}

// Prints a banner to the console
func banner() {
	fmt.Println(color.HiGreenString(`
    ____            _           ________          _     
   / __ \___ _   __(_)_______  / ____/ /_  ____ _(_)___ 
  / / / / _ \ | / / / ___/ _ \/ /   / __ \/ __  / / __ \
 / /_/ /  __/ |/ / / /__/  __/ /___/ / / / /_/ / / / / /
/_____/\___/|___/_/\___/\___/\____/_/ /_/\__,_/_/_/ /_/ 

`))
}

type Stuff struct {
	Name string
	Last string
}

func testGenerate() {
	config := core.NewDefaultInstanceConfiguration()

	raw, err := generator.GenerateInstanceConfig("my-test-config", config)
	if err != nil {
		log.Error().Err(err).Stack().Msg("unable to create configuration")
	}
	fmt.Printf("yaml content:\n\n%s", string(raw))
}
