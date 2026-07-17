// Package models its the principal data struct package to interactions, for now
package models

import (
	"github.com/bwmarrin/snowflake"
)

type Tag struct {
	TagName        string `json:"tagName"`
	TagDescription string `json:"tagDescription,omitempty"`
	TagValue       string `json:"tagValue"`
	ID             int    `json:"ID"`
}
type GroupTags struct {
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	Tags        []*Tag       `json:"tags,omitempty"`
	ID          snowflake.ID `json:"ID"`
}

type InternalUniqueValue struct {
	TagSelected *int `json:"tagSelected"`
}

func (md *InternalUniqueValue) AsValue() int {
	return *md.TagSelected
}

// func CreateTag() *Tag {
// 	node, _ := snowflake.NewNode(1)
// 	return &Tag{
// 		ID: ,
// 	}
// }

// TODO: create GroupTags
