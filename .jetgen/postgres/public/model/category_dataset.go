//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

type CategoryDataset struct {
	ID          int32 `sql:"primary_key"`
	L1In        string
	L2In        *string
	L3In        *string
	L4In        *string
	L5In        *string
	L6In        *string
	L7In        *string
	L8In        *string
	FullPathOut string
	NameOut     string
	Version     string
	Label       string
}
