package entity

import (
	"time"
)

type Post struct {
	Id       string    `json:"id"`
	Title    string    `json:"title" form:"title"`
	Body     string    `json:"body" form:"body"`
	Author   string    `json:"author" form:"author"`
	PostTime time.Time `json:"post_time,omitempty"`
}

type PostUpdate struct {
	Title  string `json:"title,omitempty" bson:"title,omitempty" form:"title"`
	Body   string `json:"body,omitempty" bson:"body,omitempty" form:"body"`
	Author string `json:"author,omitempty" bson:"author,omitempty" form:"author"`
}
