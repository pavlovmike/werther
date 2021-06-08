/*
Copyright (c) JSC iCore.

This source code is licensed under the MIT license found in the
LICENSE file in the root directory of this source tree.
*/

package hydra

import (
	"github.com/pkg/errors"
)

// ConsentReqDoer fetches information on the OAuth2 request and then accept or reject the requested authentication process.
type ConsentReqDoer struct {
	hydraURL                string
	rememberFor             int
	extendAccessTokenClaims bool
}

// NewConsentReqDoer creates a ConsentRequest.
func NewConsentReqDoer(hydraURL string, rememberFor int, extendAccessTokenClaims bool) *ConsentReqDoer {
	return &ConsentReqDoer{
		hydraURL:                hydraURL,
		rememberFor:             rememberFor,
		extendAccessTokenClaims: extendAccessTokenClaims,
	}
}

// InitiateRequest fetches information on the OAuth2 request.
func (crd *ConsentReqDoer) InitiateRequest(challenge string) (*ReqInfo, error) {
	ri, err := initiateRequest(consent, crd.hydraURL, challenge)
	return ri, errors.Wrap(err, "failed to initiate consent request")
}

// AcceptConsentRequest accepts the requested authentication process, and returns redirect URI.
func (crd *ConsentReqDoer) AcceptConsentRequest(challenge string, remember bool, grantScope []string, claims interface{},
) (string, error) {
	type session struct {
		AccessToken interface{} `json:"access_token,omitempty"`
		IDToken     interface{} `json:"id_token,omitempty"`
	}
	data := struct {
		GrantScope  []string `json:"grant_scope"`
		Remember    bool     `json:"remember"`
		RememberFor int      `json:"remember_for"`
		Session     session  `json:"session,omitempty"`
	}{
		GrantScope:  grantScope,
		Remember:    remember,
		RememberFor: crd.rememberFor,
		Session: session{
			IDToken: claims,
		},
	}
	if crd.extendAccessTokenClaims {
		data.Session.AccessToken = claims
	}
	redirectURI, err := acceptRequest(consent, crd.hydraURL, challenge, data)
	return redirectURI, errors.Wrap(err, "failed to accept consent request")
}
