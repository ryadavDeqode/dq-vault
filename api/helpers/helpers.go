package helpers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// User -- stores data related to user
type User struct {
	UUID       string `json:"uuid"`
	Mnemonic   string `json:"mnemonic"`
	Passphrase string `json:"passphrase"`
}

// CheckError checks for any potential errors
func CheckError(err error, message string) {
	if err != nil {
		log.Fatalf("%v - %v", message, err)
	}
}

// ErrMissingField returns a logical response error that prints a consistent
// error message for when a required field is missing.
func ErrMissingField(field string) *logical.Response {
	return logical.ErrorResponse(fmt.Sprintf("missing required field '%s'", field))
}

// ValidationErr returns an error that corresponds to a validation error.
func ValidationErr(msg string) error {
	return logical.CodedError(http.StatusUnprocessableEntity, msg)
}

// ValidateFields verifies that no bad arguments were given to the request.
func ValidateFields(req *logical.Request, data *framework.FieldData) error {
	var unknownFields []string
	for k := range req.Data {
		if _, ok := data.Schema[k]; !ok {
			unknownFields = append(unknownFields, k)
		}
	}

	if len(unknownFields) > 0 {
		// Sort since this is a human error
		sort.Strings(unknownFields)

		return fmt.Errorf("unknown fields: %q", unknownFields)
	}

	return nil
}

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

// New returns an error that formats as the given text.
func New(text string) error {
	return &errorString{text}
}

// ValidateUser - validates data provided by user
func ValidateUser(ctx context.Context, req *logical.Request, uuid string, derivationPath string) error {
	// Check if user provided UUID or not
	if uuid == "" {
		return errors.New("Provide a valid UUID")
	}

	// Check if user provided derivationPath or not
	if derivationPath == "" {
		return errors.New("Provide a valid path")
	}

	// Obtain all existing UUID's from DB
	vals, err := req.Storage.List(ctx, "users/")
	CheckError(err, "")

	var exists = false

	// check if UUID exists
	for i := 0; i < len(vals); i++ {
		if uuid == vals[i] {
			exists = true
			break
		}
	}

	// if UUID not exists then return error
	if !exists {
		return errors.New("Provided UUID does not exists")
	}

	return nil
}