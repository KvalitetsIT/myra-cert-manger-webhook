package models

type Record struct {

	// Action is one of 'present' or 'cleanup'.
	// If the action is 'present', the record will be presented with the
	// solving service.
	// If the action is 'cleanup', the record will be cleaned up with the
	// solving service.
	Action string `json:"action"`

	// Type is the type of ACME challenge.
	// Only dns-01 is currently supported.
	Type string `json:"type"`

	// DNSName is the name of the domain that is actually being validated, as
	// requested by the user on the Certificate resource.
	// This will be of the form 'example.com' from normal hostnames, and
	// '*.example.com' for wildcards.
	DNSName string `json:"dnsName"`

	// Key is the key that should be presented.
	// This key will already be signed by the account that owns the challenge.
	// For DNS01, this is the key that should be set for the TXT record for
	// ResolveFQDN.
	Key string `json:"key"`

	// ResourceNamespace is the namespace containing resources that are
	// referenced in the providers config.
	// If this request is solving for an Issuer resource, this will be the
	// namespace of the Issuer.
	// If this request is solving for a ClusterIssuer resource, this will be
	// the configured 'cluster resource namespace'
	ResourceNamespace string `json:"resourceNamespace"`

	// ResolvedFQDN is the fully-qualified domain name that should be
	// updated/presented after resolving all CNAMEs.
	// This should be honoured when using the DNS01 solver type.
	// This will be of the form '_acme-challenge.example.com.'.
	// +optional
	ResolvedFQDN string `json:"resolvedFQDN,omitempty"`

	// ResolvedZone is the zone encompassing the ResolvedFQDN.
	// This is included as part of the ChallengeRequest so that webhook
	// implementers do not need to implement their own SOA recursion logic.
	// This indicates the zone that the provided FQDN is encompassed within,
	// determined by performing SOA record queries for each part of the FQDN
	// until an authoritative zone is found.
	// This will be of the form 'example.com.'.
	ResolvedZone string `json:"resolvedZone,omitempty"`
}
