package dnspod

import (
	"github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/pkg/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
)

func (s *Solver) getConfigAndClient(ch *v1alpha1.ChallengeRequest) (*Config, *dnspod.Client, error) {
	cfg, err := loadConfig(ch.Config)
	if err != nil {
		s.Error(err, "failed to load config from challenge request")
		return nil, nil, errors.WithStack(err)
	}

	oidcP, _ := common.DefaultTkeOIDCRoleArnProvider()
	providers := []common.Provider{
		common.DefaultEnvProvider(),
		oidcP,
		common.DefaultProfileProvider(),
		common.DefaultCvmRoleProvider(),
	}
	provider := common.NewProviderChain(providers)
	cred, err := provider.GetCredential()
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to get credential from provider chain")
	}

	dnspodClient, err := dnspod.NewClient(cred, "", profile.NewClientProfile())
	if err != nil {
		s.Error(err, "failed to create dnspod client")
		return nil, nil, errors.WithStack(err)
	}
	return cfg, dnspodClient, nil
}
