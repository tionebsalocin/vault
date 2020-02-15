package seal

import (
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/seal"
)

func configureAliCloudKMSSeal(configKMS *configutil.KMS, infoKeys *[]string, info *map[string]string, logger log.Logger, inseal vault.Seal) (vault.Seal, error) {
	kms, kmsInfo, err := configutil.GetAliCloudKMSFunc(nil, configKMS.Config)
	if err != nil {
		return nil, err
	}
	autoseal := vault.NewAutoSeal(&seal.Access{
		Wrapper: kms,
	})
	if kmsInfo != nil {
		*infoKeys = append(*infoKeys, "Seal Type", "AliCloud KMS Region", "AliCloud KMS KeyID")
		(*info)["Seal Type"] = configKMS.Type
		(*info)["AliCloud KMS Region"] = kmsInfo["region"]
		(*info)["AliCloud KMS KeyID"] = kmsInfo["kms_key_id"]
		if domain, ok := kmsInfo["domain"]; ok {
			*infoKeys = append(*infoKeys, "AliCloud KMS Domain")
			(*info)["AliCloud KMS Domain"] = domain
		}
	}
	return autoseal, nil
}
