package search

import (
	"github.com/mises-id/mises-swapsvc/lib/db/odm"
	"github.com/mises-id/mises-swapsvc/lib/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	SwapProviderSearch struct {
		ID primitive.ObjectID
		//sort
		//limit
		ListNum int64
		//page
		Limit      int64  `json:"limit" query:"limit"`
		NextID     string `json:"last_id" query:"last_id"`
		PageNum    int64  `json:"page_num" query:"page_num"`
		PageSize   int64  `json:"page_size" query:"page_size"`
		PageParams *pagination.PageQuickParams
	}
)

func (params *SwapProviderSearch) BuildAdminSearch(chain *odm.DB) *odm.DB {
	//base
	//where
	//sort
	chain = chain.Sort(bson.M{"sort_num": -1})
	//limit
	if (params.PageNum <= 0 || params.PageSize <= 0) && params.ListNum > 0 {
		chain = chain.Limit(params.ListNum)
	}
	return chain
}

func (params *SwapProviderSearch) GetPageParams() *pagination.TraditionalParams {
	page := pagination.DefaultTraditionalParams()
	if params.PageNum > 0 {
		page.PageNum = params.PageNum
	}
	if params.PageSize > 0 {
		page.PageSize = params.PageSize
	}
	return page
}
func (params *SwapProviderSearch) GetQuickPageParams() *pagination.PageQuickParams {
	res := pagination.DefaultQuickParams()
	if params.ListNum > 0 {
		res.Limit = params.Limit
	}
	if params.NextID != "" {
		res.NextID = params.NextID
	}
	return params.PageParams
}
