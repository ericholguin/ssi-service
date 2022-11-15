package storage

import (
	"github.com/TBD54566975/ssi-sdk/credential"

	cred "github.com/tbd54566975/ssi-service/internal/credential"
	"github.com/tbd54566975/ssi-service/internal/keyaccess"
	"github.com/tbd54566975/ssi-service/internal/util"
	"github.com/tbd54566975/ssi-service/pkg/storage"
)

type StoreCredentialRequest struct {
	cred.Container
}

type StoredCredential struct {
	// This ID is generated by the storage module upon first write
	ID string `json:"id"`

	CredentialID string `json:"credentialId"`

	// only one of these fields should be present
	Credential    *credential.VerifiableCredential `json:"credential,omitempty"`
	CredentialJWT *keyaccess.JWT                   `json:"token,omitempty"`

	Issuer       string `json:"issuer"`
	Subject      string `json:"subject"`
	Schema       string `json:"schema"`
	IssuanceDate string `json:"issuanceDate"`
	Revoked      bool   `json:"revoked"`
}

func (sc StoredCredential) IsValid() bool {
	return sc.ID != "" && (sc.HasDataIntegrityCredential() || sc.HasJWTCredential())
}

func (sc StoredCredential) HasDataIntegrityCredential() bool {
	return sc.Credential != nil && sc.Credential.Proof != nil
}

func (sc StoredCredential) HasJWTCredential() bool {
	return sc.CredentialJWT != nil
}

type Storage interface {
	StoreCredential(request StoreCredentialRequest) error
	GetCredential(id string) (*StoredCredential, error)
	GetCredentialsByIssuer(issuer string) ([]StoredCredential, error)
	GetCredentialsBySubject(subject string) ([]StoredCredential, error)
	GetCredentialsBySchema(schema string) ([]StoredCredential, error)
	GetCredentialsByIssuerAndSchema(issuer, schema string) ([]StoredCredential, error)
	DeleteCredential(id string) error

	StoreStatusListCredential(request StoreCredentialRequest) error
	GetStatusListCredential(id string) (*StoredCredential, error)
	GetStatusListCredentialsByIssuerAndSchema(issuer, schema string) ([]StoredCredential, error)
	DeleteStatusListCredential(id string) error

	GetNextStatusListRandomIndex() (int, error)
}

func NewCredentialStorage(s storage.ServiceStorage) (Storage, error) {
	switch s.Type() {
	case storage.Bolt:
		gotBolt, ok := s.(*storage.BoltDB)
		if !ok {
			return nil, util.LoggingNewErrorf("trouble instantiating : %s", s.Type())
		}
		boltStorage, err := NewBoltCredentialStorage(gotBolt)
		if err != nil {
			return nil, util.LoggingErrorMsg(err, "could not instantiate credential bolt storage")
		}
		return boltStorage, err
	default:
		return nil, util.LoggingNewErrorf("unsupported storage type: %s", s.Type())
	}
}
