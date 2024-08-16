package validations

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	domainUser "github.com/trio-kwek-kwek/GoWhatsappWeb/domains/user"
	pkgError "github.com/trio-kwek-kwek/GoWhatsappWeb/pkg/error"
)

func ValidateUserInfo(ctx context.Context, request domainUser.InfoRequest) error {
	err := validation.ValidateStructWithContext(ctx, &request,
		validation.Field(&request.Phone, validation.Required),
	)

	if err != nil {
		return pkgError.ValidationError(err.Error())
	}

	return nil
}
func ValidateUserAvatar(ctx context.Context, request domainUser.AvatarRequest) error {
	err := validation.ValidateStructWithContext(ctx, &request,
		validation.Field(&request.Phone, validation.Required),
		validation.Field(&request.IsCommunity, validation.When(request.IsCommunity, validation.Required, validation.In(true, false))),
		validation.Field(&request.IsPreview, validation.When(request.IsPreview, validation.Required, validation.In(true, false))),
	)

	if err != nil {
		return pkgError.ValidationError(err.Error())
	}

	return nil
}
