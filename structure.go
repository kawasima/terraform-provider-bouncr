package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)
func expandStringSet(configured *schema.Set) []string {
	return expandStringList(configured.List())
}

func expandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, val)
		}
	}
	return vs
}
