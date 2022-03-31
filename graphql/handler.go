/*
Copyright © 2022 SiteWhere LLC - All Rights Reserved
Unauthorized copying of this file, via any medium is strictly prohibited.
Proprietary and confidential.
*/

package graphql

import (
	"context"
	"net/http"

	graphql "github.com/graph-gophers/graphql-go"

	"github.com/graph-gophers/graphql-go/relay"
)

type ContextKey string

const (
	ContextRdbKey ContextKey = "rdb"
)

// Adds extra context to http request.
type HttpHandler struct {
	Schema           *graphql.Schema
	Relay            *relay.Handler
	ContextProviders map[ContextKey]interface{}
}

// Create new http handler.
func NewHttpHandler(schema *graphql.Schema, providers map[ContextKey]interface{}) *HttpHandler {
	handler := &HttpHandler{
		Schema:           schema,
		Relay:            &relay.Handler{Schema: schema},
		ContextProviders: providers,
	}
	return handler
}

// Handles http request processing.
func (h *HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for key, value := range h.ContextProviders {
		r = r.WithContext(context.WithValue(r.Context(), key, value))
	}
	h.Relay.ServeHTTP(w, r)
}
