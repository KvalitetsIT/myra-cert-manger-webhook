package adaptors

import (
	"github.com/KvalitetsIT/cert-manager-webhook-core/pkg/client"
	"github.com/KvalitetsIT/cert-manager-webhook-core/pkg/client/adaptors"
	"github.com/KvalitetsIT/cert-manager-webhook-core/pkg/models"
	"github.com/KvalitetsIT/cert-manager-webhook-myra/internal/mapping"
	"github.com/Myra-Security-GmbH/myrasec-go/v2"
)

type MyraClientAdaptor struct {
	a *adaptors.ClientAdaptor[models.Record, myrasec.DNSRecord]
}

func NewMyraClientAdaptor(client client.Client[myrasec.DNSRecord]) *MyraClientAdaptor {
	return &MyraClientAdaptor{
		adaptors.NewClientAdaptor(
			client,
			mapping.MyraMapper{},
		),
	}
}

func (p *MyraClientAdaptor) OnDelete(record models.Record) (models.Record, error) {
	return p.a.OnDelete(record)
}

func (p *MyraClientAdaptor) OnAdd(record models.Record) (models.Record, error) {
	return p.a.OnAdd(record)
}
