package service

import (
	"os"

	"github.com/SennaSemakula/nrobservgen/pkg/hcl"
)

// NewTemplate generates a template
func NewYamlTemplate(service, dest string) (string, error) {
	// Move this logic out
	svc := NewServiceDefaults(service)
	out, err := WriteYaml(svc, dest)
	if err != nil {
		return "", err
	}
	return out, nil
}

func NewTFTemplate(svc Service, dest string) (string, error) {
	if err := hcl.WriteVars(svc.Name, svc.Version, dest); err != nil {
		return "", err
	}
	if err := hcl.WriteProvider(svc.Name, svc.Version, dest); err != nil {
		return "", err
	}
	if err := hcl.WriteBackend(svc.Name, svc.Version, dest); err != nil {
		return "", err
	}

	// TODO: Make this more abstract.
	// These are also hardcoded sources. Should change this
	//  Have defaultvars for alerts and dashboard modules
	// TODO: uncomment once custom alerts are supported in terraform modules.
	// Depedent on what Mike Price is working on.
	vars := hcl.DefaultVars(svc.Name)
	localVars := hcl.DefaultLocalVars(svc.Name)
	modules := []hcl.Module{
		{
			Name:      "alerts",
			Source:    "git::ssh://github.com/<terraform_module>",
			Version:   svc.Version,
			Variables: vars,
			LocalVars: localVars,
		},
		{
			Name:      "dashboards",
			Source:    "git::ssh://github.com/<terraform_module>",
			Version:   svc.Version,
			Variables: vars,
		},
	}
	if err := hcl.WriteMain(svc.Name, svc.Version, dest, modules); err != nil {
		return "", err
	}

	if len(dest) == 0 {
		dir, err := os.Getwd()
		if err != nil {
			return "", err
		}
		return dir, nil
	}

	return dest, nil
}
