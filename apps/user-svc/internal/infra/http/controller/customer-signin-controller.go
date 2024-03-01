package controller

import (
	"net/http"
	"user-svc/internal/application/usecases"
	"user-svc/internal/domain/customer"
	"user-svc/internal/infra/adapters"
	"user-svc/internal/infra/database"
	"user-svc/internal/infra/http/helper"
)

func CustomerSignin(w http.ResponseWriter, r *http.Request) {
	body, err := helper.ParseRequestBody[customer.SigninCustomerIn](r)
	if err != nil {
		helper.HandleHttpError(w, err)
		return
	}

	err = validateSigninPayload(body)
	if err != nil {
		helper.HandleHttpError(w, err)
		return
	}

	repo := database.NewPgxCustomerRepository(database.Conn)
	hasher := adapters.NewBcryptPasswordHasher()
	tkGen := adapters.NewJwtTokenGenerator()
	customerSigninService := usecases.NewCustomerSigninService(repo, hasher, tkGen)

	res, err := customerSigninService.Execute(body)
	if err != nil {
		helper.HandleHttpError(w, err)
		return
	}

	helper.HandleHttpSuccessJson(w, http.StatusOK, res)
}

func validateSigninPayload(in *customer.SigninCustomerIn) error {
	if in.Email == "" || in.Password == "" {
		return helper.ErrInvalidArgument
	}
	return nil
}
