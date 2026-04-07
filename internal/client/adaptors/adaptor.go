package adaptors

import (
	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/client"
	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/mapping"
)

// The purpose of this struct is to adapt the myra specific api into the internal models

type clientAdaptor[internal, external any] struct {
	client client.Client[external]
	mapper mapping.Mapper[internal, external]
}

func (p *clientAdaptor[Internal, External]) OnDelete(record Internal) (Internal, error) {
	r, e := p.client.OnDelete(p.mapper.ToExternal(record))
	return p.mapper.ToInternal(r), e
}

func (p *clientAdaptor[Internal, External]) OnAdd(record Internal) (Internal, error) {
	r, e := p.client.OnAdd(p.mapper.ToExternal(record))
	return p.mapper.ToInternal(r), e
}
