package jsonprovider

import (
	"strings"

	"github.com/hashicorp/terraform/configs/configschema"
)

var catalogMap = map[string]string{
	"ECS":   "Compute",
	"AS":    "Compute",
	"EVS":   "Storage",
	"OBS":   "Storage",
	"SFS":   "Storage",
	"VPC":   "Network",
	"ELB":   "Network",
	"DNS":   "Network",
	"NAT":   "Network",
	"VPCEP": "Network",
	"DDS":   "Database",
	"RDS":   "Database",
	"CCE":   "Container",
	"APIG":  "Application",
}

var serviceMap = map[string]string{
	"ECS":   "Elastic Cloud Server",
	"AS":    "Auto Scaling",
	"EVS":   "Elastic Volume Service",
	"OBS":   "Object Storage Service",
	"SFS":   "Scalable File Service",
	"VPC":   "Virtual Private Cloud",
	"ELB":   "Elastic Load Balance",
	"DNS":   "Domain Name Service",
	"NAT":   "NAT Gateway",
	"VPCEP": "VPC Endpoint",
	"DDS":   "Document Database Service",
	"RDS":   "Relational Database Service",
	"CCE":   "Cloud Container Engine",
	"APIG":  "API GateWay",
}

var resourceMap = map[string]string{
	"compute":             "ECS",
	"networking":          "VPC",
	"networking_vip":      "VPC",
	"networking_secgroup": "VPC",
	"vpc_peering":         "VPC",
	"lb":                  "ELB",
	"rds_read_replica":    "RDS",
	"api_gateway":         "APIG",
	"nat_dnat":            "NAT",
	"nat_snat":            "NAT",
	"obs_bucket":          "OBS",
	"sfs_access":          "SFS",
	"sfs_file":            "SFS",
}

type schema struct {
	Version uint64   `json:"version,omitempty"`
	Block   *block   `json:"block,omitempty"`
	Product *product `json:"product,omitempty"`
}

type product struct {
	Catalog string `json:"catalog,omitempty"`
	Name    string `json:"name,omitempty"`
	Short   string `json:"short,omitempty"`
}

// marshalSchema is a convenience wrapper around mashalBlock. Schema version
// should be set by the caller.
func marshalSchema(block *configschema.Block) *schema {
	if block == nil {
		return &schema{}
	}

	var ret schema
	ret.Block = marshalBlock(block)

	return &ret
}

func marshalSchemas(blocks map[string]*configschema.Block, rVersions map[string]uint64) map[string]*schema {
	if blocks == nil {
		return map[string]*schema{}
	}
	ret := make(map[string]*schema, len(blocks))
	for k, v := range blocks {
		ret[k] = marshalSchema(v)
		version, ok := rVersions[k]
		if ok {
			ret[k].Version = version
		}
	}
	return ret
}

func marshalSchemaByName(blocks map[string]*configschema.Block, rVersions map[string]uint64, name string) map[string]*schema {
	if blocks == nil {
		return map[string]*schema{}
	}
	ret := make(map[string]*schema, len(blocks))
	for k, v := range blocks {
		if k == name {
			ret[k] = marshalSchema(v)
			version, ok := rVersions[k]
			if ok {
				ret[k].Version = version
			}
			product := marshalProduct(k)
			if product != nil {
				ret[k].Product = product
			}
		}
	}
	return ret
}

func marshalProduct(key string) *product {
	var name string
	nameSlice := strings.Split(key, "_")
	switch len(nameSlice) {
	case 1:
		return nil
	case 2, 3:
		name = nameSlice[1]
	default:
		name = strings.Join(nameSlice[1:len(nameSlice)-1], "_")
	}

	shortName := resourceMap[name]
	if shortName == "" {
		shortName = strings.ToUpper(name)
	}
	if serviceMap[shortName] == "" || catalogMap[shortName] == "" {
		return nil
	}

	return &product{
		Short:   shortName,
		Name:    serviceMap[shortName],
		Catalog: catalogMap[shortName],
	}
}
