// generated by stringer -type FeatureCompat; DO NOT EDIT

package main

import "fmt"

const _FeatureCompat_name = "MaintainsChecksums"

var _FeatureCompat_index = [...]uint8{0, 18}

func (i FeatureCompat) String() string {
	i -= 1
	if i+1 >= FeatureCompat(len(_FeatureCompat_index)) {
		return fmt.Sprintf("FeatureCompat(%d)", i+1)
	}
	return _FeatureCompat_name[_FeatureCompat_index[i]:_FeatureCompat_index[i+1]]
}
