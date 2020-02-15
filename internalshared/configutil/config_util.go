// +build !enterprise

package configutil

import (
	"github.com/hashicorp/errwrap"
	wrapping "github.com/hashicorp/go-kms-wrapping"
	"github.com/hashicorp/go-kms-wrapping/wrappers/alicloudkms"
	"github.com/hashicorp/go-kms-wrapping/wrappers/awskms"
	"github.com/hashicorp/go-kms-wrapping/wrappers/azurekeyvault"
	"github.com/hashicorp/go-kms-wrapping/wrappers/gcpckms"
	"github.com/hashicorp/go-kms-wrapping/wrappers/ocikms"
	"github.com/hashicorp/go-kms-wrapping/wrappers/transit"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/vault/sdk/logical"
)

type EntSharedConfig struct {
}

func (ec *EntSharedConfig) ParseConfig(list *ast.ObjectList) error {
	return nil
}

func ParseEntropy(result *SharedConfig, list *ast.ObjectList, blockName string) error {
	return nil
}

var GetAWSKMSFunc = func(opts *wrapping.WrapperOptions, config map[string]string) (wrapping.Wrapper, map[string]string, error) {
	kms := awskms.NewWrapper(nil)
	kmsInfo, err := kms.SetConfig(config)
	if err != nil {
		// If the error is any other than logical.KeyNotFoundError, return the error
		if !errwrap.ContainsType(err, new(logical.KeyNotFoundError)) {
			return nil, nil, err
		}
	}
	return kms, kmsInfo, nil
}

var GetAliCloudKMSFunc = func(opts *wrapping.WrapperOptions, config map[string]string) (wrapping.Wrapper, map[string]string, error) {
	kms := alicloudkms.NewWrapper(nil)
	kmsInfo, err := kms.SetConfig(config)
	if err != nil {
		// If the error is any other than logical.KeyNotFoundError, return the error
		if !errwrap.ContainsType(err, new(logical.KeyNotFoundError)) {
			return nil, nil, err
		}
	}
	return kms, kmsInfo, nil
}

var GetAzureKMSFunc = func(opts *wrapping.WrapperOptions, config map[string]string) (wrapping.Wrapper, map[string]string, error) {
	kms := azurekeyvault.NewWrapper(nil)
	kmsInfo, err := kms.SetConfig(config)
	if err != nil {
		// If the error is any other than logical.KeyNotFoundError, return the error
		if !errwrap.ContainsType(err, new(logical.KeyNotFoundError)) {
			return nil, nil, err
		}
	}
	return kms, kmsInfo, nil
}

var GetGCPCKMSFunc = func(opts *wrapping.WrapperOptions, config map[string]string) (wrapping.Wrapper, map[string]string, error) {
	kms := gcpckms.NewWrapper(nil)
	kmsInfo, err := kms.SetConfig(config)
	if err != nil {
		// If the error is any other than logical.KeyNotFoundError, return the error
		if !errwrap.ContainsType(err, new(logical.KeyNotFoundError)) {
			return nil, nil, err
		}
	}
	return kms, kmsInfo, nil
}

var GetOCIKMSFunc = func(opts *wrapping.WrapperOptions, config map[string]string) (wrapping.Wrapper, map[string]string, error) {
	kms := ocikms.NewWrapper(nil)
	kmsInfo, err := kms.SetConfig(config)
	if err != nil {
		return nil, nil, err
	}
	return kms, kmsInfo, nil
}

var GetTransitKMSFunc = func(opts *wrapping.WrapperOptions, config map[string]string) (wrapping.Wrapper, map[string]string, error) {
	kms := transit.NewWrapper(opts)
	kmsInfo, err := kms.SetConfig(config)
	if err != nil {
		// If the error is any other than logical.KeyNotFoundError, return the error
		if !errwrap.ContainsType(err, new(logical.KeyNotFoundError)) {
			return nil, nil, err
		}
	}
	return kms, kmsInfo, nil
}
