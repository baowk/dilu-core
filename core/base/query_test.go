package base

import (
	"fmt"
	"testing"
)

type TP struct {
	Id       int    `json:"id" query:""`
	Name     string `json:"name" query:"type:like;"`
	Status   int    `json:"status" query:"type:gt"`
	DeptPath string `json:"deptPath" query:"type:left;"`
	File     string `json:"file" query:"type:right;"`
	Flag     int    `json:"flag" query:"type:gt"`
	Uid      []int  `json:"uid" query:"type:in"`
}

func (TP) TableName() string {
	return "tp"
}

func TestResolveSearchQuery(t *testing.T) {
	tp := TP{
		Id:       1,
		Name:     "abc",
		Status:   1,
		DeptPath: "/a/b/",
		File:     ".png",
	}
	condition := &GormCondition{
		GormPublic: GormPublic{},
		Join:       make([]*GormJoin, 0),
	}
	ResolveSearchQuery("mysql", tp, condition)
	for _, join := range condition.Join {
		if join == nil {
			continue
		}
		fmt.Println(join.JoinOn)
		for k, v := range join.Where {
			fmt.Println(k, v)
		}
		for k, v := range join.Or {
			fmt.Println(k, v)
		}
		for _, o := range join.Order {
			fmt.Println(o)
		}
	}
	for k, v := range condition.Where {
		fmt.Println(k, v)
	}
	for k, v := range condition.Or {
		fmt.Println(k, v)
	}
	for _, o := range condition.Order {
		fmt.Println(o)
	}
}
