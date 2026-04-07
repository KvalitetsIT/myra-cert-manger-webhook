package adaptors

import (
	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/client"
	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/mapping"
	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/models"
	"github.com/Myra-Security-GmbH/myrasec-go/v2"
)

type MyraClientAdaptor struct {
	*clientAdaptor[models.Record, myrasec.DNSRecord]
}

func NewMyraClientAdaptor(client client.Client[myrasec.DNSRecord]) *MyraClientAdaptor {
	return &MyraClientAdaptor{
		&clientAdaptor[models.Record, myrasec.DNSRecord]{
			client: client,
			mapper: &mapping.MyraMapper{},
		},
	}
}

func (p *MyraClientAdaptor) OnDelete(record models.Record) (models.Record, error) {
	var r myrasec.DNSRecord = p.mapper.ToExternal(record)
	r, err := p.client.OnDelete(r)
	return p.mapper.ToInternal(r), err
}

func (p *MyraClientAdaptor) OnAdd(record models.Record) (models.Record, error) {
	r, e := p.client.OnAdd(p.mapper.ToExternal(record))
	return p.mapper.ToInternal(r), e
}
