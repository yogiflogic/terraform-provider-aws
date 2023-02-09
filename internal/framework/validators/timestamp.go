package validators

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type utcTimestampValidator struct{}

func (validator utcTimestampValidator) Description(_ context.Context) string {
	return "value must be a valid UTC Timestamp"
}

func (validator utcTimestampValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

func (validator utcTimestampValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	if err := validateUTCTimestamp(request.ConfigValue.ValueString()); err != nil {
		response.Diagnostics.Append(diag.NewAttributeErrorDiagnostic(
			request.Path,
			validator.Description(ctx),
			err.Error(),
		))
		return
	}
}

func UTCTimestamp() validator.String {
	return utcTimestampValidator{}
}

func validateUTCTimestamp(value string) error {
	_, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return fmt.Errorf("must be in RFC3339 time format %q. Example: %s", time.RFC3339, err)
	}

	return nil
}

type onceADayWindowFormatValidator struct{}

func (validator onceADayWindowFormatValidator) Description(_ context.Context) string {
	return "value must be a valid time format"
}

func (validator onceADayWindowFormatValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

func (validator onceADayWindowFormatValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	if err := validateOnceADayWindowFormat(request.ConfigValue.ValueString()); err != nil {
		response.Diagnostics.Append(diag.NewAttributeErrorDiagnostic(
			request.Path,
			validator.Description(ctx),
			err.Error(),
		))
		return
	}
}

func OnceADayWindowFormat() validator.String {
	return onceADayWindowFormatValidator{}
}

func validateOnceADayWindowFormat(value string) error {
	// valid time format is "hh24:mi"
	validTimeFormat := "([0-1][0-9]|2[0-3]):([0-5][0-9])"
	validTimeFormatConsolidated := "^(" + validTimeFormat + "-" + validTimeFormat + "|)$"

	if !regexp.MustCompile(validTimeFormatConsolidated).MatchString(value) {
		return fmt.Errorf("(%s) must satisfy the format of \"hh24:mi-hh24:mi\"", value)
	}

	return nil
}

type onceAWeekWindowFormatValidator struct{}

func (validator onceAWeekWindowFormatValidator) Description(_ context.Context) string {
	return "value must be a valid time format"
}

func (validator onceAWeekWindowFormatValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

func (validator onceAWeekWindowFormatValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	if err := validateOnceAWeekWindowFormat(request.ConfigValue.ValueString()); err != nil {
		response.Diagnostics.Append(diag.NewAttributeErrorDiagnostic(
			request.Path,
			validator.Description(ctx),
			err.Error(),
		))
		return
	}
}

func OnceAWeekWindowFormat() validator.String {
	return onceAWeekWindowFormatValidator{}
}

func validateOnceAWeekWindowFormat(value string) error {
	// valid time format is "ddd:hh24:mi"
	validTimeFormat := "(sun|mon|tue|wed|thu|fri|sat):([0-1][0-9]|2[0-3]):([0-5][0-9])"
	validTimeFormatConsolidated := "^(" + validTimeFormat + "-" + validTimeFormat + "|)$"

	val := strings.ToLower(value)
	if !regexp.MustCompile(validTimeFormatConsolidated).MatchString(val) {
		return fmt.Errorf("(%s) must satisfy the format of \"ddd:hh24:mi-ddd:hh24:mi\"", value)
	}
	return nil
}
