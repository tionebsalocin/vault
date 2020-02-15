package seal

import (
	"fmt"

	"github.com/hashicorp/go-hclog"
	wrapping "github.com/hashicorp/go-kms-wrapping"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/vault"
)

var (
	ConfigureSeal = configureSeal
)

func configureSeal(configKMS *configutil.KMS, infoKeys *[]string, info *map[string]string, logger hclog.Logger, inseal vault.Seal) (outseal vault.Seal, err error) {
	switch configKMS.Type {
	case wrapping.AliCloudKMS:
		return configureAliCloudKMSSeal(configKMS, infoKeys, info, logger, inseal)

	case wrapping.AWSKMS:
		return configureAWSKMSSeal(configKMS, infoKeys, info, logger, inseal)

	case wrapping.AzureKeyVault:
		return configureAzureKeyVaultSeal(configKMS, infoKeys, info, logger, inseal)

	case wrapping.GCPCKMS:
		return configureGCPCKMSSeal(configKMS, infoKeys, info, logger, inseal)

	case wrapping.OCIKMS:
		return configureOCIKMSSeal(configKMS, infoKeys, info, logger, inseal)

	case wrapping.Transit:
		return configureTransitSeal(configKMS, infoKeys, info, logger, inseal)

	case wrapping.PKCS11:
		return nil, fmt.Errorf("Seal type 'pkcs11' requires the Vault Enterprise HSM binary")

	case wrapping.Shamir:
		return inseal, nil

	default:
		return nil, fmt.Errorf("Unknown seal type %q", configKMS.Type)
	}
}
