/*
Copyright © 2022 SiteWhere LLC - All Rights Reserved
Unauthorized copying of this file, via any medium is strictly prohibited.
Proprietary and confidential.
*/

package core

import (
	"fmt"
	"os"

	"github.com/devicechain-io/dc-k8s/api/v1beta1"
	"github.com/rs/zerolog/log"
)

// Gets the k8s tenantmicroservice resource based on values passed in the environment
func (ms *Microservice) getTenantMicroservice() (*v1beta1.TenantMicroservice, error) {
	instanceId, exists := os.LookupEnv(ENV_INSTANCE_ID)
	if !exists {
		return nil, fmt.Errorf("instance id not passed via environment variable: Looking in %s", ENV_INSTANCE_ID)
	}
	tmid, exists := os.LookupEnv(ENV_TENANTMICROSERVICE_ID)
	if !exists {
		return nil, fmt.Errorf("tenant microservice id not passed via environment variable: Looking in %s", ENV_INSTANCE_ID)
	}

	tm, err := v1beta1.GetTenantMicroservice(v1beta1.TenantMicroserviceGetRequest{
		InstanceId:           instanceId,
		TenantMicroserviceId: tmid,
	})
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("Unable to locate tenantmicroservice referenced in environment variable: %s", tmid))
		return nil, err
	}
	return tm, nil
}
