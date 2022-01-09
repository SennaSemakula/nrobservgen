package hcl

import (
	"fmt"

	"github.com/SennaSemakula/nrobservgen/pkg/common"
)

func WriteMain(service, version, dest string, modules []Module) error {
	root := NewBody(version)

	for _, v := range modules {
		root.WriteModule(v)
	}

	// Ugly. Refactor
	name := "main.tf"
	if len(dest) > 0 {
		if _, err := common.WriteFile(service, dest+"/"+name, root.file.Bytes()); err != nil {
			return err
		}
		return nil
	}

	if _, err := common.WriteFile(service, name, root.file.Bytes()); err != nil {
		return err
	}
	return nil

}

func WriteVars(service, version, dest string) error {
	// TODO: stop passing in service
	root := NewBody(version)
	vars := DefaultVars(service)

	for _, v := range vars {
		if err := root.WriteVar(v); err != nil {
			return err
		}
	}

	// Ugly. Refactor
	name := "variables.tf"
	if len(dest) > 0 {
		if _, err := common.WriteFile(service, dest+"/"+name, root.file.Bytes()); err != nil {
			return err
		}
		return nil
	}

	if _, err := common.WriteFile(service, name, root.file.Bytes()); err != nil {
		return err
	}
	return nil
}

func WriteProvider(service, version, dest string) error {
	root := NewBody(version)
	p := DefaultProviders()

	tfBlock := root.AppendNewBlock("terraform", nil).Body()
	providerBlock := tfBlock.AppendNewBlock("required_providers", nil)

	b := Block{providerBlock}
	for _, v := range p {
		if err := b.WriteProvider(v); err != nil {
			return err
		}
	}

	// Ugly. Refactor.
	// TODO: stop passing in service
	name := "provider.tf"
	if len(dest) > 0 {
		common.WriteFile(service, dest+"/"+name, root.file.Bytes())
		return nil
	}

	common.WriteFile(service, name, root.file.Bytes())
	return nil
}

func WriteBackend(service, version, dest string) error {
	// TODO: stop passing in service
	root := NewBody(version)

	tfBlock := root.AppendNewBlock("terraform", nil).Body()
	// modularise this. currently hardcoded to just s3
	tfBlock.AppendNewBlock("backend", []string{"s3"})

	// Ugly. Refactor
	name := "backend.tf"
	if len(dest) > 0 {
		if _, err := common.WriteFile(service, dest+"/"+name, root.file.Bytes()); err != nil {
			return err
		}
		return nil
	}

	if _, err := common.WriteFile(service, name, root.file.Bytes()); err != nil {
		return err
	}
	return nil
}

func DefaultProviders() []*Provider {
	// TODO: pin providers
	return []*Provider{
		{"newrelic", "newrelic/newrelic"},
		{"aws", "hashicorp/aws"},
	}
}

func DefaultVars(service string) []*Variable {
	return []*Variable{
		{"newrelic_account_id", "number", 0},
		{"newrelic_region", "string", "EU"},
		{"newrelic_environment", "string", ""},
		{"newrelic_api_key", "string", ""},
		{"aws_region", "string", "eu-west-1"},
		{"newrelic_local_creds", "bool", false},
		{"service_settings_file", "string", fmt.Sprintf("%s-alerts.yaml", service)},
	}
}

func DefaultLocalVars(service string) []*Variable {
	return []*Variable{
		{"service_name", "string", "local.settings.service_name"}, // why is the terraform module requring service_name local var???.
		{"nrql_alerts", "string", "try(local.settings.alerts, {})"},
	}
}

func NewModule(name, source, version string, vars []*Variable) Module {
	return Module{
		Name:      name,
		Source:    source,
		Version:   version,
		Variables: vars,
	}
}
