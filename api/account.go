package api

import (
	"database/sql"
	"net/http"

	db "github.com/ferseg/golang-simple-bank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type CreateAccountRequest struct {
	Owner string `json:"owner" binding:"required"`
	// Balance  int64  `json:"balance" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, arg)

	if err != nil {
    if pqError, ok := err.(*pq.Error); ok {
      switch pqError.Code.Name() {
      case "foreign_key_violation", "unique_violation":
        ctx.JSON(http.StatusForbidden, error(err))
        return
      }
    }
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type GetAccountRequest struct {
  ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
  var req GetAccountRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

  account, err := server.store.GetAccount(ctx, req.ID)

  if err != nil {
    if err == sql.ErrNoRows {
      ctx.JSON(http.StatusNotFound, "Account not found")
      return
    }
    ctx.JSON(http.StatusInternalServerError, errorResponse(err))
    return
  }

  ctx.JSON(http.StatusOK, account)
}

type ListAccountsRequest struct {
	PageNumber int32 `form:"page_num" binding:"required,min=1"`
	PageSize   int32 `form:"page_size" binding:"required,min=5,max=15"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var req ListAccountsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageSize - 1) * req.PageNumber,
	}

	account, err := server.store.ListAccounts(ctx, arg)
  if err != nil {
    ctx.JSON(http.StatusBadRequest, errorResponse(err))
    return
  }

  ctx.JSON(http.StatusOK, account)
}
