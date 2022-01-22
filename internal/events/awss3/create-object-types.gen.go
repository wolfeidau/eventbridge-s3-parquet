// Package awss3 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.1 DO NOT EDIT.
package awss3

import (
	"time"
)

// AWSEvent defines model for AWSEvent.
type AWSEvent struct {
	Account    string        `json:"account"`
	Detail     ObjectCreated `json:"detail"`
	DetailType string        `json:"detail-type"`
	Id         string        `json:"id"`
	Region     string        `json:"region"`
	Resources  []string      `json:"resources"`
	Source     string        `json:"source"`
	Time       time.Time     `json:"time"`
	Version    string        `json:"version"`
}

// Bucket defines model for Bucket.
type Bucket struct {
	Name string `json:"name"`
}

// Object defines model for Object.
type Object struct {
	Etag      string  `json:"etag"`
	Key       string  `json:"key"`
	Sequencer string  `json:"sequencer"`
	Size      float32 `json:"size"`
	VersionId *string `json:"version-id,omitempty"`
}

// ObjectCreated defines model for ObjectCreated.
type ObjectCreated struct {
	Bucket          Bucket `json:"bucket"`
	Object          Object `json:"object"`
	Reason          string `json:"reason"`
	RequestId       string `json:"request-id"`
	Requester       string `json:"requester"`
	SourceIpAddress string `json:"source-ip-address"`
	Version         string `json:"version"`
}
