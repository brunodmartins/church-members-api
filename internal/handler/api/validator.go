package api

import (
	"fmt"
	apierrors "github.com/BrunoDM2943/church-members-api/platform/infra/errors"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
)

func Validate(dto interface{}) error {
	validate := validator.New()
	err := validate.Struct(dto)
	if err != nil {
		builder := strings.Builder{}
		for _, err := range err.(validator.ValidationErrors) {
			builder.WriteString(fmt.Sprintf("Field:%s, Tag:%s,Value:%s\n", err.StructNamespace(), err.Tag(), err.Param()))
		}
		return apierrors.NewApiError(builder.String(), http.StatusBadRequest)
	}
	return nil
}
