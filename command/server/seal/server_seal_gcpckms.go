package seal

import (
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/sdk/helper/useragent"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/seal"
)

func configureGCPCKMSSeal(configKMS *configutil.KMS, infoKeys *[]string, info *map[string]string, logger log.Logger, inseal vault.Seal) (vault.Seal, error) {
	configKMS.Config["user_agent"] = useragent.String()
	kms, kmsInfo, err := configutil.GetGCPCKMSFunc(nil, configKMS.Config)
	if err != nil {
		return nil, err
	}
	autoseal := vault.NewAutoSeal(&seal.Access{
		Wrapper: kms,
	})
	if kmsInfo != nil {
		*infoKeys = append(*infoKeys, "Seal Type", "GCP KMS Project", "GCP KMS Region", "GCP KMS Key Ring", "GCP KMS Crypto Key")
		(*info)["Seal Type"] = configKMS.Type
		(*info)["GCP KMS Project"] = kmsInfo["project"]
		(*info)["GCP KMS Region"] = kmsInfo["region"]
		(*info)["GCP KMS Key Ring"] = kmsInfo["key_ring"]
		(*info)["GCP KMS Crypto Key"] = kmsInfo["crypto_key"]
	}
	return autoseal, nil
}
