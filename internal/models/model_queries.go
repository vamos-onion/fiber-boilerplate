package models

import (
	"context"
	"errors"

	"gopkg.in/guregu/null.v4"
)

func query() *Queries {
	return New(SQL)
}

var Appuser AppuserQuery = new(AppuserBlock)

var ArrayTest ArrayTestQuery = new(ArrayTestBlock)

type AppuserQuery interface {
	GetAppusersByName(qctx context.Context, exid string) (AppuserBlock, error)
	SearchAppusers(qctx context.Context, param SearchAppusersParams) ([]AppuserBlock, error)
	CreateAppuser(tx *Queries, qctx context.Context, param CreateAppuserParams) (AppuserBlock, error)
	UpdateAppuser(tx *Queries, qctx context.Context, param UpdateAppuserParams) (AppuserBlock, error)
}

type ArrayTestQuery interface {
	GetAllColumns(qctx context.Context) ([]ArrayTestBlock, error)
}

// AppuserBlock :
func (m *AppuserBlock) GetAppusersByName(qctx context.Context, exid string) (AppuserBlock, error) {
	if qctx == nil {
		qctx = context.Background()
	}
	return query().GetAppusersByName(qctx, null.StringFrom(exid))
}

func (m *AppuserBlock) SearchAppusers(qctx context.Context, param SearchAppusersParams) ([]AppuserBlock, error) {
	if qctx == nil {
		qctx = context.Background()
	}
	return query().SearchAppusers(qctx, param)
}

func (m *AppuserBlock) CreateAppuser(tx *Queries, qctx context.Context, param CreateAppuserParams) (AppuserBlock, error) {
	if tx == nil {
		if qctx == nil {
			qctx = context.Background()
		}
		return query().CreateAppuser(qctx, param)
	} else {
		if qctx == nil {
			return AppuserBlock{}, errors.New("qctx is nil")
		}
		return tx.CreateAppuser(qctx, param)
	}
}

func (m *AppuserBlock) UpdateAppuser(tx *Queries, qctx context.Context, param UpdateAppuserParams) (AppuserBlock, error) {
	if tx == nil {
		if qctx == nil {
			qctx = context.Background()
		}
		return query().UpdateAppuser(qctx, param)
	} else {
		if qctx == nil {
			return AppuserBlock{}, errors.New("qctx is nil")
		}
		return tx.UpdateAppuser(qctx, param)
	}
}

// ArrayTestBlock :
func (m *ArrayTestBlock) GetAllColumns(qctx context.Context) ([]ArrayTestBlock, error) {
	if qctx == nil {
		qctx = context.Background()
	}
	return query().GetAllColumns(qctx)
}
