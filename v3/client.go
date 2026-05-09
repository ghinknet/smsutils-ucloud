package ucloud

import (
	"github.com/ucloud/ucloud-sdk-go/services/usms"
	"github.com/ucloud/ucloud-sdk-go/ucloud/auth"
	"github.com/ucloud/ucloud-sdk-go/ucloud/config"
	"github.com/ucloud/ucloud-sdk-go/ucloud/log"
	"go.gh.ink/smsutils/v3/errors"
	"go.gh.ink/smsutils/v3/model"
)

type Client struct {
	Client *usms.USMSClient
	// JSON
	Marshal   func(any) ([]byte, error)
	Unmarshal func([]byte, any) error
}

type Driver struct{}

func (d Driver) NewClient(params model.DriverClientParam) (model.Client, error) {
	// Check credential
	pubKey, priKey, prjID := params.Credential[PublicKey], params.Credential[PrivateKey], params.Credential[ProjectId]
	if pubKey == "" || priKey == "" || prjID == "" {
		return Client{}, errors.ErrDriverCredentialInvalid.WithDriverName(Name)
	}

	cfg := config.NewConfig()
	cfg.LogLevel = log.FatalLevel
	cfg.ProjectId = prjID

	// Set private key and public key
	credential := auth.NewCredential()
	credential.PublicKey = pubKey
	credential.PrivateKey = priKey

	return Client{
		Client:    usms.NewClient(&cfg, &credential),
		Marshal:   params.Marshal,
		Unmarshal: params.Unmarshal,
	}, nil
}
