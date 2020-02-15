package seal

import (
	log "github.com/hashicorp/go-hclog"
	wrapping "github.com/hashicorp/go-kms-wrapping"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/seal"
)

func configureTransitSeal(configKMS *configutil.KMS, infoKeys *[]string, info *map[string]string, logger log.Logger, inseal vault.Seal) (vault.Seal, error) {
	kms, kmsInfo, err := configutil.GetTransitKMSFunc(
		&wrapping.WrapperOptions{
			Logger: logger.ResetNamed("seal-transit"),
		}, configKMS.Config)
	if err != nil {
		return nil, err
	}
	autoseal := vault.NewAutoSeal(&seal.Access{
		Wrapper: kms,
	})
	if kmsInfo != nil {
		*infoKeys = append(*infoKeys, "Seal Type", "Transit Address", "Transit Mount Path", "Transit Key Name")
		(*info)["Seal Type"] = configKMS.Type
		(*info)["Transit Address"] = kmsInfo["address"]
		(*info)["Transit Mount Path"] = kmsInfo["mount_path"]
		(*info)["Transit Key Name"] = kmsInfo["key_name"]
		if namespace, ok := kmsInfo["namespace"]; ok {
			*infoKeys = append(*infoKeys, "Transit Namespace")
			(*info)["Transit Namespace"] = namespace
		}
	}
	return autoseal, nil
}
