package utils

import (
	// "fmt"
	// "strconv"
	// "strings"

	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
)

// BuildReference create reference from defined object
func BuildReference(uuid, kind string) *v3.Reference {
	return &v3.Reference{
		Kind: StringPtr(kind),
		UUID: StringPtr(uuid),
	}
}
