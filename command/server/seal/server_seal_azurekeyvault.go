package seal

import (
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/seal"
)

func configureAzureKeyVaultSeal(configKMS *configutil.KMS, infoKeys *[]string, info *map[string]string, logger log.Logger, inseal vault.Seal) (vault.Seal, error) {
	kms, kmsInfo, err := configutil.GetAzureKMSFunc(nil, configKMS.Config)
	if err != nil {
		return nil, err
	}
	autoseal := vault.NewAutoSeal(&seal.Access{
		Wrapper: kms,
	})
	if kmsInfo != nil {
		*infoKeys = append(*infoKeys, "Seal Type", "Azure Environment", "Azure Vault Name", "Azure Key Name")
		(*info)["Seal Type"] = configKMS.Type
		(*info)["Azure Environment"] = kmsInfo["environment"]
		(*info)["Azure Vault Name"] = kmsInfo["vault_name"]
		(*info)["Azure Key Name"] = kmsInfo["key_name"]
	}
	return autoseal, nil
}
