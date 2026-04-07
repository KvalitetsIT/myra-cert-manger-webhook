package mapping_test

import (
	"testing"

	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/mapping"
	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/models"
	"github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/stretchr/testify/require"
)

func TestCertManagerMapper_ToInternal(t *testing.T) {
	mapper := mapping.NewCertManagerMapper()

	tests := []struct {
		name     string
		input    v1alpha1.ChallengeRequest
		expected models.Record
	}{
		{
			name: "present action with type",
			input: v1alpha1.ChallengeRequest{
				Action:            v1alpha1.ChallengeActionPresent,
				Type:              "dns-01",
				DNSName:           "example.com",
				Key:               "key123",
				ResourceNamespace: "default",
				ResolvedFQDN:      "_acme-challenge.example.com.",
				ResolvedZone:      "example.com.",
			},
			expected: models.Record{
				Action:            "present",
				Type:              "dns-01",
				DNSName:           "example.com",
				Key:               "key123",
				ResourceNamespace: "default",
				ResolvedFQDN:      "_acme-challenge.example.com.",
				ResolvedZone:      "example.com.",
			},
		},
		{
			name: "cleanup action with empty type defaults to TXT",
			input: v1alpha1.ChallengeRequest{
				Action:            v1alpha1.ChallengeActionCleanUp,
				Type:              "",
				DNSName:           "example.org",
				Key:               "key456",
				ResourceNamespace: "default",
				ResolvedFQDN:      "_acme-challenge.example.org.",
				ResolvedZone:      "example.org.",
			},
			expected: models.Record{
				Action:            "cleanup",
				Type:              "TXT",
				DNSName:           "example.org",
				Key:               "key456",
				ResourceNamespace: "default",
				ResolvedFQDN:      "_acme-challenge.example.org.",
				ResolvedZone:      "example.org.",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.ToInternal(tt.input)
			require.Equal(t, tt.expected.Action, result.Action, "Action field mismatch")
			require.Equal(t, tt.expected.Type, result.Type, "Type field mismatch")
			require.Equal(t, tt.expected.DNSName, result.DNSName, "DNSName field mismatch")
			require.Equal(t, tt.expected.Key, result.Key, "Key field mismatch")
			require.Equal(t, tt.expected.ResourceNamespace, result.ResourceNamespace, "ResourceNamespace field mismatch")
			require.Equal(t, tt.expected.ResolvedFQDN, result.ResolvedFQDN, "ResolvedFQDN field mismatch")
			require.Equal(t, tt.expected.ResolvedZone, result.ResolvedZone, "ResolvedZone field mismatch")
		})
	}
}

func TestCertManagerMapper_ToExternal(t *testing.T) {
	mapper := mapping.NewCertManagerMapper()

	tests := []struct {
		name     string
		input    models.Record
		expected v1alpha1.ChallengeRequest
	}{
		{
			name: "present action",
			input: models.Record{
				Action:            "present",
				Type:              "dns-01",
				DNSName:           "example.com",
				Key:               "key123",
				ResourceNamespace: "default",
				ResolvedFQDN:      "_acme-challenge.example.com.",
				ResolvedZone:      "example.com.",
			},
			expected: v1alpha1.ChallengeRequest{
				Action:            v1alpha1.ChallengeActionPresent,
				Type:              "dns-01",
				DNSName:           "example.com",
				Key:               "key123",
				ResourceNamespace: "default",
				ResolvedFQDN:      "_acme-challenge.example.com.",
				ResolvedZone:      "example.com.",
			},
		},
		{
			name: "cleanup action",
			input: models.Record{
				Action:            "cleanup",
				Type:              "TXT",
				DNSName:           "example.org",
				Key:               "key456",
				ResourceNamespace: "default",
				ResolvedFQDN:      "_acme-challenge.example.org.",
				ResolvedZone:      "example.org.",
			},
			expected: v1alpha1.ChallengeRequest{
				Action:            v1alpha1.ChallengeActionCleanUp,
				Type:              "TXT",
				DNSName:           "example.org",
				Key:               "key456",
				ResourceNamespace: "default",
				ResolvedFQDN:      "_acme-challenge.example.org.",
				ResolvedZone:      "example.org.",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.ToExternal(tt.input)
			require.Equal(t, tt.expected.Action, result.Action, "Action field mismatch")
			require.Equal(t, tt.expected.Type, result.Type, "Type field mismatch")
			require.Equal(t, tt.expected.DNSName, result.DNSName, "DNSName field mismatch")
			require.Equal(t, tt.expected.Key, result.Key, "Key field mismatch")
			require.Equal(t, tt.expected.ResourceNamespace, result.ResourceNamespace, "ResourceNamespace field mismatch")
			require.Equal(t, tt.expected.ResolvedFQDN, result.ResolvedFQDN, "ResolvedFQDN field mismatch")
			require.Equal(t, tt.expected.ResolvedZone, result.ResolvedZone, "ResolvedZone field mismatch")
		})
	}
}
