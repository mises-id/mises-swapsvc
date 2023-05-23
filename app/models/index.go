package models

import (
	"github.com/mises-id/mises-swapsvc/lib/db/odm"
	"github.com/mises-id/mises-swapsvc/lib/pagination"
)

type (
	IAdminParams interface {
		BuildAdminSearch(chain *odm.DB) *odm.DB
	}
	IAdminPageParams interface {
		BuildAdminSearch(chain *odm.DB) *odm.DB
		GetPageParams() *pagination.TraditionalParams
	}
	IAdminQuickPageParams interface {
		BuildAdminSearch(chain *odm.DB) *odm.DB
		GetQuickPageParams() *pagination.PageQuickParams
	}
)
