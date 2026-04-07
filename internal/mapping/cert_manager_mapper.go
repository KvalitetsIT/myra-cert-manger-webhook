package mapping

import (
	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/models"
	"github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
)

type CertManagerMapper struct{}

func NewCertManagerMapper() *CertManagerMapper {
	return &CertManagerMapper{}
}

func (c *CertManagerMapper) ToInternal(ch v1alpha1.ChallengeRequest) models.Record {
	action := map[v1alpha1.ChallengeAction]string{
		v1alpha1.ChallengeActionPresent: "present",
		v1alpha1.ChallengeActionCleanUp: "cleanup",
	}[ch.Action]

	return models.Record{
		Action:            action,
		Type:              get_record_type(ch.Type),
		DNSName:           ch.DNSName,
		Key:               ch.Key,
		ResourceNamespace: ch.ResourceNamespace,
		ResolvedFQDN:      ch.ResolvedFQDN,
		ResolvedZone:      ch.ResolvedZone,
	}
}

// Since the myra SDK expects a record type during creation this is added.
// Assuming the record type is included in the request from cert manager this has no effect.
// However since it ain't included in the integration test (./cmd/integration/integration_test.go) this is required.
// See: https://github.com/cert-manager/webhook-example/issues/3
func get_record_type(s string) string {
	if s != "" {
		return s
	}
	return "TXT"
}

func (c *CertManagerMapper) ToExternal(ch models.Record) v1alpha1.ChallengeRequest {
	action := map[string]v1alpha1.ChallengeAction{
		"present": v1alpha1.ChallengeActionPresent,
		"cleanup": v1alpha1.ChallengeActionCleanUp,
	}[ch.Action]

	return v1alpha1.ChallengeRequest{
		Action:            action,
		Type:              ch.Type,
		DNSName:           ch.DNSName,
		Key:               ch.Key,
		ResourceNamespace: ch.ResourceNamespace,
		ResolvedFQDN:      ch.ResolvedFQDN,
		ResolvedZone:      ch.ResolvedZone,
	}
}
