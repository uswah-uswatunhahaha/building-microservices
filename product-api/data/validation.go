package data

import (
	"regexp"

	"github.com/go-playground/validator"
)

// Validate to validate incoming value
func (d *Product) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("sku", skuValidation)
	return validate.Struct(d)
}

func skuValidation(fl validator.FieldLevel) bool {

	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := re.FindAllString(fl.Field().String(), -1)
	if len(matches) != 1 {
		return false
	}

	return true
}
