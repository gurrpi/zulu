package api

import (
	"errors"

	"github.com/DE-labtory/zulu/account"
	"github.com/DE-labtory/zulu/keychain"
	"github.com/DE-labtory/zulu/types"
	"github.com/gin-gonic/gin"
)

type WalletApi struct {
	resolver  *account.Resolver
	generator keychain.KeyGenerator
	store     keychain.Store
}

func NewWallet(resolver *account.Resolver) *WalletApi {
	return &WalletApi{
		resolver: resolver,
	}
}

func (w *WalletApi) CreateWallet(context *gin.Context) {
	key, err := w.generator.Generate()
	if err != nil {
		internalServerError(context, errors.New("failed to generate key pair"))
	}
	accounts, err := w.getAccounts(key)
	if err != nil {
		internalServerError(context, errors.New("failed to derive account"))
		return
	}

	if err := w.store.Store(key); err != nil {
		internalServerError(context, errors.New("failed to store key"))
		return
	}
	context.JSON(200, accounts)
}

func (w *WalletApi) GetWallet(context *gin.Context) {
	var request struct {
		ID string `uri:"id" binding:"required"`
	}

	if err := context.ShouldBindUri(&request); err != nil {
		badrequestError(context, errors.New("path variable :id does not exists"))
		return
	}

	key, err := w.store.Get(request.ID)
	if err != nil {
		internalServerError(context, errors.New("failed to read key"))
		return
	}

	accounts, err := w.getAccounts(key)
	if err != nil {
		internalServerError(context, errors.New("failed to derive account"))
		return
	}
	context.JSON(200, accounts)
}

func (w *WalletApi) getAccounts(key keychain.Key) ([]types.Account, error) {
	var accounts []types.Account
	for _, s := range w.resolver.GetAllServices() {
		a, err := s.DeriveAccount(key)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}
	return accounts, nil
}

func (w *WalletApi) Transfer(context *gin.Context) {
	context.JSON(200, gin.H{
		"message": "transfer",
	})
}
