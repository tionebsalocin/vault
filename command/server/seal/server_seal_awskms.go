package seal

import (
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/seal"
)

func configureAWSKMSSeal(configKMS *configutil.KMS, infoKeys *[]string, info *map[string]string, logger hclog.Logger, inseal vault.Seal) (vault.Seal, error) {
	kms, kmsInfo, err := configutil.GetAWSKMSFunc(nil, configKMS.Config)
	if err != nil {
		return nil, err
	}
	autoseal := vault.NewAutoSeal(&seal.Access{
		Wrapper: kms,
	})
	if kmsInfo != nil {
		*infoKeys = append(*infoKeys, "Seal Type", "AWS KMS Region", "AWS KMS KeyID")
		(*info)["Seal Type"] = configKMS.Type
		(*info)["AWS KMS Region"] = kmsInfo["region"]
		(*info)["AWS KMS KeyID"] = kmsInfo["kms_key_id"]
		if endpoint, ok := kmsInfo["endpoint"]; ok {
			*infoKeys = append(*infoKeys, "AWS KMS Endpoint")
			(*info)["AWS KMS Endpoint"] = endpoint
		}
	}
	return autoseal, nil
}
