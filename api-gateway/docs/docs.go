// Package docs Spotify Clone API
//
// REST API dokumentacija za Spotify Clone aplikaciju.
//
//	Schemes: http, https
//	Host: localhost:8080
//	BasePath: /api/v1
//	Version: 1.0.0
//
//	SecurityDefinitions:
//	  BearerAuth:
//	    type: apiKey
//	    name: Authorization
//	    in: header
//
// swagger:meta
package docs

import "embed"

//go:embed swagger.json
var SwaggerJSON embed.FS
