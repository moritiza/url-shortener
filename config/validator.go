package config

import (
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

// Validator create new validator with en language
func Validator() *validator.Validate {
	v := validator.New()
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(v, trans)

	return v
}

// ValidatorErrors return validator errors and translate them to en language
func ValidatorErrors(v *validator.Validate, err error) string {
	var errString string

	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	validatorCustomErrors(v, trans)

	// Get all validator errors
	errs := err.(validator.ValidationErrors)

	// Translate each validator error
	for _, e := range errs {
		errString += e.Translate(trans) + "\n"
	}

	return strings.TrimSuffix(errString, "\n")
}

// validatorCustomErrors set custom error message for validate tags
func validatorCustomErrors(v *validator.Validate, trans ut.Translator) {
	v.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field!", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())

		return t
	})

	v.RegisterTranslation("url", trans, func(ut ut.Translator) error {
		return ut.Add("url", "{0} not a valid URL!", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("url", fe.Field())

		return t
	})
}
