// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package errors

import (
	"github.com/berachain/beacon-kit/mod/errors"
)

var (
	// ErrPreDefinedJSONRPC is a catch-all error for all pre-defined json-rpc
	// errors.
	ErrPreDefinedJSONRPC = errors.New(
		"json-rpc error",
	)

	// ErrUnknownPayload indicates an unavailable or non-existent payload
	// (JSON-RPC code -38001).
	ErrUnknownPayload = errors.New(
		"payload does not exist or is not available")

	// ErrInvalidForkchoiceState indicates an invalid fork choice state
	// (JSON-RPC code -38002).
	ErrInvalidForkchoiceState = errors.New(
		"invalid forkchoice state")

	// ErrInvalidPayloadAttributes indicates invalid or inconsistent payload
	// attributes
	// (JSON-RPC code -38003).
	ErrInvalidPayloadAttributes = errors.New(
		"payload attributes are invalid / inconsistent")

	// ErrRequestTooLarge indicates that the request is too large
	// (JSON-RPC code -38004).
	ErrRequestTooLarge = errors.New(
		"request is too large",
	)

	// ErrUnknownPayloadStatus indicates an unknown payload status.
	ErrUnknownPayloadStatus = errors.New(
		"unknown payload status")

	// ErrAcceptedPayloadStatus indicates a payload status of ACCEPTED.
	ErrAcceptedPayloadStatus = errors.New(
		"payload status is ACCEPTED")

	// ErrSyncingPayloadStatus indicates a payload status of SYNCING.
	ErrSyncingPayloadStatus = errors.New(
		"payload status is SYNCING",
	)

	// ErrInvalidPayloadStatus indicates an invalid payload status.
	ErrInvalidPayloadStatus = errors.New(
		"payload status is INVALID")

	// ErrInvalidBlockHashPayloadStatus indicates a failure in validating the
	// block hash for the payload.
	ErrInvalidBlockHashPayloadStatus = errors.New(
		"payload status is INVALID_BLOCK_HASH")

	// ErrNilForkchoiceResponse indicates a nil forkchoice response.
	ErrNilForkchoiceResponse = errors.New(
		"nil forkchoice response",
	)
	/// ErrNilBlobsBundle is returned when nil blobs bundle is received.
	ErrNilBlobsBundle = errors.New(
		"nil blobs bundle received from execution client")

	// ErrInvalidPayloadAttributeVersion indicates an invalid version of payload
	// attributes was provided.
	ErrInvalidPayloadAttributeVersion = errors.New(
		"invalid payload attribute version")
	// ErrInvalidPayloadType indicates an invalid payload type
	// was provided for an RPC call.
	ErrInvalidPayloadType = errors.New("invalid payload type for RPC call")

	// ErrInvalidGetPayloadVersion indicates that an unknown fork version was
	// provided for getting a payload.
	ErrInvalidGetPayloadVersion = errors.New("unknown fork for get payload")

	// ErrUnsupportedVersion indicates a request for a block type with an
	// unknown ExecutionPayload schema.
	ErrUnsupportedVersion = errors.New(
		"unknown ExecutionPayload schema for block version")

	// ErrNilJWTSecret indicates that the JWT secret is nil.
	ErrNilJWTSecret = errors.New("nil JWT secret")

	// ErrNilAttributesPassedToClient is returned when nil attributes are
	// passed to the client.
	ErrNilAttributesPassedToClient = errors.New(
		"nil attributes passed to client",
	)

	// ErrNilExecutionPayloadEnvelope is returned when nil execution payload
	// envelope is received.
	ErrNilExecutionPayloadEnvelope = errors.New(
		"nil execution payload envelope received from execution client")

	// ErrNilExecutionPayload is returned when nil execution payload
	// envelope is received.
	ErrNilExecutionPayload = errors.New(
		"nil execution payload received from execution client")

	// ErrNilPayloadStatus is returned when nil payload status is received.
	ErrNilPayloadStatus = errors.New(
		"nil payload status received from execution client",
	)

	// ErrEngineAPITimeout is returned when the engine API call times out.
	ErrEngineAPITimeout = errors.New(
		"engine API call timed out",
	)
)
