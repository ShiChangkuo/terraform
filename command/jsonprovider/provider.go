package jsonprovider

import (
	"encoding/json"

	"github.com/hashicorp/terraform/terraform"
)

const (
	TypeProvider int = iota
	TypeResource
	TypeData
)

// providers is the top-level object returned when exporting provider schemas
type providers struct {
	Schemas map[string]*Provider `json:"provider_schemas,omitempty"`
}

type Provider struct {
	Provider          *schema            `json:"provider,omitempty"`
	ResourceSchemas   map[string]*schema `json:"resource_schemas,omitempty"`
	DataSourceSchemas map[string]*schema `json:"data_source_schemas,omitempty"`
}

func newProviders() *providers {
	schemas := make(map[string]*Provider)
	return &providers{
		Schemas: schemas,
	}
}

func Marshal(s *terraform.Schemas, sourceType int, name string) ([]byte, error) {
	providers := newProviders()

	for k, v := range s.Providers {
		switch sourceType {
		case TypeResource:
			providers.Schemas[k.String()] = marshalResource(v, name)
		case TypeData:
			providers.Schemas[k.String()] = marshalDataSource(v, name)
		default:
			providers.Schemas[k.String()] = marshalProvider(v)
		}

	}

	ret, err := json.MarshalIndent(providers, "", "  ")
	return ret, err
}

func marshalProvider(tps *terraform.ProviderSchema) *Provider {
	if tps == nil {
		return &Provider{}
	}

	var ps *schema
	if tps.Provider != nil {
		ps = marshalSchema(tps.Provider)
	}

	return &Provider{
		Provider: ps,
	}
}

func marshalResource(tps *terraform.ProviderSchema, key string) *Provider {
	if tps == nil {
		return &Provider{}
	}

	var rs map[string]*schema

	if tps.ResourceTypes != nil {
		if key == "all" {
			rs = marshalSchemas(tps.ResourceTypes, tps.ResourceTypeSchemaVersions)
		} else {
			rs = marshalSchemaByName(tps.ResourceTypes, tps.ResourceTypeSchemaVersions, key)
		}
	}

	return &Provider{
		ResourceSchemas: rs,
	}
}

func marshalDataSource(tps *terraform.ProviderSchema, key string) *Provider {
	if tps == nil {
		return &Provider{}
	}

	var ds map[string]*schema
	if tps.DataSources != nil {
		if key == "all" {
			ds = marshalSchemas(tps.DataSources, tps.ResourceTypeSchemaVersions)
		} else {
			ds = marshalSchemaByName(tps.DataSources, tps.ResourceTypeSchemaVersions, key)
		}
	}

	return &Provider{
		DataSourceSchemas: ds,
	}
}
