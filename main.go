/*
Copyright © 2022 SiteWhere LLC - All Rights Reserved
Unauthorized copying of this file, via any medium is strictly prohibited.
Proprietary and confidential.
*/

package main

import (
	"context"
	"os"

	"github.com/devicechain-io/dc-microservice/core"
	"github.com/rs/zerolog/log"
)

// Create and run a microservice
func main() {
	os.Setenv(core.ENV_INSTANCE_ID, "dc1")
	os.Setenv(core.ENV_TENANTMICROSERVICE_ID, "tms-tenant1-devicemanagement")
	callbacks := core.LifecycleCallbacks{
		Initializer: core.LifecycleCallback{
			Preprocess: func(ctx context.Context) error {
				log.Info().Msg("Microservice called pre-initialize callback.")
				return nil
			},
			Postprocess: func(ctx context.Context) error {
				log.Info().Msg("Microservice called post-initialize callback.")
				return nil
			},
		},
		Starter: core.LifecycleCallback{
			Preprocess: func(ctx context.Context) error {
				log.Info().Msg("Microservice called pre-start callback.")
				return nil
			},
			Postprocess: func(ctx context.Context) error {
				log.Info().Msg("Microservice called post-start callback.")
				return nil
			},
		},
		Stopper: core.LifecycleCallback{
			Preprocess: func(ctx context.Context) error {
				log.Info().Msg("Microservice called pre-stop callback.")
				return nil
			},
			Postprocess: func(ctx context.Context) error {
				log.Info().Msg("Microservice called post-stop callback.")
				return nil
			},
		},
		Terminator: core.LifecycleCallback{
			Preprocess: func(ctx context.Context) error {
				log.Info().Msg("Microservice called pre-terminate callback.")
				return nil
			},
			Postprocess: func(ctx context.Context) error {
				log.Info().Msg("Microservice called post-terminate callback.")
				return nil
			},
		},
	}
	ms := core.NewMicroservice("test-service", callbacks)
	ms.Run()
}
