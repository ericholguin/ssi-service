package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tbd54566975/ssi-service/pkg/service/operation/storage"
)

var didWebContext = NewTestContext("DIDWeb")

func TestCreateIssuerDIDWebIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	didWebOutput, err := CreateDIDWeb()
	assert.NoError(t, err)

	issuerDID, err := getJSONElement(didWebOutput, "$.did.id")
	assert.NoError(t, err)
	assert.Contains(t, issuerDID, "did:web")
	SetValue(didWebContext, "issuerDID", issuerDID)

	issuerKID, err := getJSONElement(didWebOutput, "$.did.verificationMethod[0].id")
	assert.NoError(t, err)
	assert.NotEmpty(t, issuerKID)
	SetValue(didWebContext, "issuerKID", issuerKID)
}

func TestCreateAliceDIDWebIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	didWebOutput, err := CreateDIDWeb()
	assert.NoError(t, err)
	assert.NotEmpty(t, didWebOutput)

	aliceDID, err := getJSONElement(didWebOutput, "$.did.id")
	assert.NoError(t, err)
	assert.Contains(t, aliceDID, "did:web")
	SetValue(didWebContext, "aliceDID", aliceDID)

	aliceKID, err := getJSONElement(didWebOutput, "$.did.verificationMethod[0].id")
	assert.NoError(t, err)
	assert.NotEmpty(t, aliceKID)
	SetValue(didWebContext, "aliceKID", aliceKID)

	aliceDIDPrivateKey, err := getJSONElement(didWebOutput, "$.privateKeyBase58")
	assert.NoError(t, err)
	assert.NotEmpty(t, aliceDID)
	SetValue(didWebContext, "aliceDIDPrivateKey", aliceDIDPrivateKey)
}

func TestDIDWebCreateSchemaIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	output, err := CreateKYCSchema()
	assert.NoError(t, err)

	schemaID, err := getJSONElement(output, "$.id")
	assert.NoError(t, err)
	assert.NotEmpty(t, schemaID)
	SetValue(didWebContext, "schemaID", schemaID)
}

func TestDIDWebCreateVerifiableCredentialIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	issuerDID, err := GetValue(didWebContext, "issuerDID")
	assert.NoError(t, err)
	assert.NotEmpty(t, issuerDID)

	issuerKID, err := GetValue(didWebContext, "issuerKID")
	assert.NoError(t, err)
	assert.NotEmpty(t, issuerKID)

	schemaID, err := GetValue(didWebContext, "schemaID")
	assert.NoError(t, err)
	assert.NotEmpty(t, schemaID)

	vcOutput, err := CreateVerifiableCredential(credInputParams{
		IssuerID:  issuerDID.(string),
		IssuerKID: issuerKID.(string),
		SchemaID:  schemaID.(string),
		SubjectID: issuerDID.(string),
	}, false)
	assert.NoError(t, err)
	assert.NotEmpty(t, vcOutput)

	credentialJWT, err := getJSONElement(vcOutput, "$.credentialJwt")
	assert.NoError(t, err)
	assert.NotEmpty(t, credentialJWT)
	SetValue(didWebContext, "credentialJWT", credentialJWT)
}

func TestDIDWebCreateCredentialManifestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	issuerDID, err := GetValue(didWebContext, "issuerDID")
	assert.NoError(t, err)
	assert.NotEmpty(t, issuerDID)

	issuerKID, err := GetValue(didWebContext, "issuerKID")
	assert.NoError(t, err)
	assert.NotEmpty(t, issuerKID)

	schemaID, err := GetValue(didWebContext, "schemaID")
	assert.NoError(t, err)
	assert.NotEmpty(t, schemaID)

	cmOutput, err := CreateCredentialManifest(credManifestParams{
		IssuerID:  issuerDID.(string),
		IssuerKID: issuerKID.(string),
		SchemaID:  schemaID.(string),
	})
	assert.NoError(t, err)

	presentationDefinitionID, err := getJSONElement(cmOutput, "$.credential_manifest.presentation_definition.id")
	assert.NoError(t, err)
	assert.NotEmpty(t, presentationDefinitionID)
	SetValue(didWebContext, "presentationDefinitionID", presentationDefinitionID)

	manifestID, err := getJSONElement(cmOutput, "$.credential_manifest.id")
	assert.NoError(t, err)
	assert.NotEmpty(t, manifestID)
	SetValue(didWebContext, "manifestID", manifestID)
}

func TestDIDWebSubmitAndReviewApplicationIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	credentialJWT, err := GetValue(didWebContext, "credentialJWT")
	assert.NoError(t, err)
	assert.NotEmpty(t, credentialJWT)

	presentationDefinitionID, err := GetValue(didWebContext, "presentationDefinitionID")
	assert.NoError(t, err)
	assert.NotEmpty(t, presentationDefinitionID)

	manifestID, err := GetValue(didWebContext, "manifestID")
	assert.NoError(t, err)
	assert.NotEmpty(t, manifestID)

	aliceDID, err := GetValue(didWebContext, "aliceDID")
	assert.NoError(t, err)
	assert.NotEmpty(t, aliceDID)

	aliceKID, err := GetValue(didWebContext, "aliceKID")
	assert.NoError(t, err)
	assert.NotEmpty(t, aliceKID)

	aliceDIDPrivateKey, err := GetValue(didWebContext, "aliceDIDPrivateKey")
	assert.NoError(t, err)
	assert.NotEmpty(t, aliceDIDPrivateKey)

	credAppJWT, err := CreateCredentialApplicationJWT(credApplicationParams{
		DefinitionID: presentationDefinitionID.(string),
		ManifestID:   manifestID.(string),
	}, credentialJWT.(string), aliceDID.(string), aliceKID.(string), aliceDIDPrivateKey.(string))
	assert.NoError(t, err)
	assert.NotEmpty(t, credAppJWT)

	submitApplicationOutput, err := SubmitApplication(applicationParams{
		ApplicationJWT: credAppJWT,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, submitApplicationOutput)

	isDone, err := getJSONElement(submitApplicationOutput, "$.done")
	assert.NoError(t, err)
	assert.Equal(t, "false", isDone)
	opID, err := getJSONElement(submitApplicationOutput, "$.id")
	assert.NoError(t, err)

	reviewApplicationOutput, err := ReviewApplication(reviewApplicationParams{
		ID:       storage.StatusObjectID(opID),
		Approved: true,
		Reason:   "oh yeah im testing",
	})
	assert.NoError(t, err)
	crManifestID, err := getJSONElement(reviewApplicationOutput, "$.credential_response.manifest_id")
	assert.NoError(t, err)
	assert.Equal(t, manifestID, crManifestID)

	vc, err := getJSONElement(reviewApplicationOutput, "$.verifiableCredentials[0]")
	assert.NoError(t, err)
	assert.NotEmpty(t, vc)

	operationOutput, err := get(endpoint + version + "operations/" + opID)
	assert.NoError(t, err)
	isDone, err = getJSONElement(operationOutput, "$.done")
	assert.NoError(t, err)
	assert.Equal(t, "true", isDone)

	opCredentialResponse, err := getJSONElement(operationOutput, "$.result.response")
	assert.NoError(t, err)
	assert.JSONEq(t, reviewApplicationOutput, opCredentialResponse)
}
