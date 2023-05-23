package search

import (
	"strings"

	"github.com/mises-id/mises-swapsvc/app/models/enum"
	"github.com/mises-id/mises-swapsvc/lib/db/odm"
	"github.com/mises-id/mises-swapsvc/lib/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	SwapOrderSearch struct {
		ID            primitive.ObjectID
		ChainID       uint64
		ReceiptStatus enum.SwapOrderStatus
		Hash          string
		FromAddress   string
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

func (params *SwapOrderSearch) BuildAdminSearch(chain *odm.DB) *odm.DB {
	//base
	//where
	if params.ChainID > 0 {
		chain = chain.Where(bson.M{"chain_id": params.ChainID})
	}
	if params.ReceiptStatus > 0 {
		chain = chain.Where(bson.M{"status": params.ReceiptStatus})
	}
	if params.Hash != "" {
		chain = chain.Where(bson.M{"tx.hash": params.Hash})
	}
	if params.FromAddress != "" {
		chain = chain.Where(bson.M{"from_address": strings.ToLower(params.FromAddress)})
	}
	//sort
	chain = chain.Sort(bson.M{"_id": -1})
	//limit
	if (params.PageNum <= 0 || params.PageSize <= 0) && params.ListNum > 0 {
		chain = chain.Limit(params.ListNum)
	}
	return chain
}

func (params *SwapOrderSearch) GetPageParams() *pagination.TraditionalParams {
	page := pagination.DefaultTraditionalParams()
	if params.PageNum > 0 {
		page.PageNum = params.PageNum
	}
	if params.PageSize > 0 {
		page.PageSize = params.PageSize
	}
	return page
}
func (params *SwapOrderSearch) GetQuickPageParams() *pagination.PageQuickParams {
	res := pagination.DefaultQuickParams()
	if params.ListNum > 0 {
		res.Limit = params.Limit
	}
	if params.NextID != "" {
		res.NextID = params.NextID
	}
	return params.PageParams
}
