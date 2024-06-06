package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strings"

	en "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	validator "gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

// validate is used to validate structs during HTTP request decoding.
var validate = validator.New()

// translator is initialized with the English locale as the default language.
var translator *ut.UniversalTranslator

// init initializes the global validator and translator.
func init() {
	// Create a new English locale.
	enLocale := en.New()

	// Create a new universal translator with English as the default language.
	translator = ut.New(enLocale, enLocale)

	// Get the translator for the English language.
	lang, _ := translator.GetTranslator("en")

	// Register the default translations for the validator.
	// The validator uses the translator to translate the error messages.
	en_translations.RegisterDefaultTranslations(validate, lang)

	// Register a custom tag name function for the validator.
	// The function extracts the JSON field name from the struct field tag.
	// If the field is not tagged with "json", the field is ignored.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		// Split the field tag by comma and get the first part.
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		// If the field is not tagged with "json", return an empty string to ignore it.
		if name == "-" {
			return ""
		}

		// Return the JSON field name.
		return name
	})
}

// Decode decodes the request body into the given struct.
// It also validates the struct and returns an error if the validation fails.
// The error returned by this function is of type *web.Error and contains the
// validation error messages.
func Decode(r *http.Request, val interface{}) error {
	// Create a new JSON decoder from the request body
	decoder := json.NewDecoder(r.Body)

	// Disallow unknown fields in the decoded struct
	decoder.DisallowUnknownFields()

	// Decode the request body into the given struct
	if err := decoder.Decode(&val); err != nil {
		// If decoding fails, return a *web.RequestError with the appropriate status code
		return NewRequestError(err, http.StatusBadRequest)
	}

	// Validate the decoded struct
	if err := validate.Struct(val); err != nil {
		// If the validation fails, get the validation errors
		verrors, ok := err.(validator.ValidationErrors)

		// If the error is not of type ValidationErrors, return the error
		if !ok {
			return err
		}

		// Get the English language translator
		lang, _ := translator.GetTranslator("en")

		// Initialize an empty slice to store the validation error messages
		var fields []FieldError

		// Iterate over the validation errors
		for _, v := range verrors {
			// Get the JSON field name from the struct tag
			field := FieldError{
				Field: v.Field(),
				Error: v.Translate(lang),
			}

			// Append the error message to the slice
			fields = append(fields, field)
		}

		// Return a *web.Error with the appropriate status code and the validation error messages
		return &Error{
			Err:    errors.New("field validation error"),
			Status: http.StatusBadRequest,
			Fields: fields,
		}
	}

	// Return nil if the validation is successful
	return nil
}
