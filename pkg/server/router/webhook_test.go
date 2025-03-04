package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tbd54566975/ssi-service/config"
	"github.com/tbd54566975/ssi-service/pkg/service/framework"
	"github.com/tbd54566975/ssi-service/pkg/service/webhook"
)

func TestWebhookRouter(t *testing.T) {
	t.Run("Nil Service", func(tt *testing.T) {
		webhookRouter, err := NewWebhookRouter(nil)
		assert.Error(tt, err)
		assert.Empty(tt, webhookRouter)
		assert.Contains(tt, err.Error(), "service cannot be nil")
	})

	t.Run("Bad Service", func(tt *testing.T) {
		webhookRouter, err := NewWebhookRouter(&testService{})
		assert.Error(tt, err)
		assert.Empty(tt, webhookRouter)
		assert.Contains(tt, err.Error(), "could not create webhook router with service type: test")
	})

	t.Run("Webhook Service Test", func(tt *testing.T) {

		db := setupTestDB(tt)
		assert.NotNil(tt, db)

		serviceConfig := config.WebhookServiceConfig{}
		webhookService, err := webhook.NewWebhookService(serviceConfig, db)
		assert.NoError(tt, err)
		assert.NotEmpty(tt, webhookService)

		// check type and status
		assert.Equal(tt, framework.Webhook, webhookService.Type())
		assert.Equal(tt, framework.StatusReady, webhookService.Status().Status)
	})

}
