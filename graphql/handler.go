/*
Copyright © 2022 SiteWhere LLC - All Rights Reserved
Unauthorized copying of this file, via any medium is strictly prohibited.
Proprietary and confidential.
*/

package graphql

import (
	"context"
	"net/http"

	"github.com/devicechain-io/dc-microservice/rdb"
	graphql "github.com/graph-gophers/graphql-go"

	"github.com/graph-gophers/graphql-go/relay"
)

type ContextKey string

const (
	ContextRdbKey ContextKey = "rdb"
)

// Adds extra context to http request.
type HttpHandler struct {
	Schema     *graphql.Schema
	Relay      *relay.Handler
	RdbManager *rdb.RdbManager
}

// Create new http handler.
func NewHttpHandler(schema *graphql.Schema, rdbmgr *rdb.RdbManager) *HttpHandler {
	handler := &HttpHandler{
		Schema:     schema,
		Relay:      &relay.Handler{Schema: schema},
		RdbManager: rdbmgr,
	}
	return handler
}

// Handles http request processing.
func (h *HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r = r.WithContext(context.WithValue(r.Context(), ContextRdbKey, h.RdbManager))
	h.Relay.ServeHTTP(w, r)
}
