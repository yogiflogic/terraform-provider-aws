package validators

import (
	"context"
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type arnValidator struct{}

func (validator arnValidator) Description(_ context.Context) string {
	return "value must be a valid ARN"
}

func (validator arnValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

func (validator arnValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	if err := validateARN(request.ConfigValue.ValueString()); err != nil {
		response.Diagnostics.Append(diag.NewAttributeErrorDiagnostic(
			request.Path,
			validator.Description(ctx),
			err.Error(),
		))
		return
	}
}

func ARN() validator.String {
	return arnValidator{}
}

var (
	accountIDRegexp = regexp.MustCompile(`^(aws|aws-managed|\d{12})$`)
	partitionRegexp = regexp.MustCompile(`^aws(-[a-z]+)*$`)
	regionRegexp    = regexp.MustCompile(`^[a-z]{2}(-[a-z]+)+-\d$`)
)

func validateARN(value string) error {
	if value == "" {
		return nil
	}

	parsedARN, err := arn.Parse(value)

	if err != nil {
		return fmt.Errorf("(%s) is an invalid ARN: %s", value, err)
	}

	if parsedARN.Partition == "" {
		return fmt.Errorf("(%s) is an invalid ARN: missing partition value", value)
	} else if !partitionRegexp.MatchString(parsedARN.Partition) {
		return fmt.Errorf("(%s) is an invalid ARN: invalid partition value (expecting to match regular expression: %s)", value, partitionRegexp)
	}

	if parsedARN.Region != "" && !regionRegexp.MatchString(parsedARN.Region) {
		return fmt.Errorf("(%s) is an invalid ARN: invalid region value (expecting to match regular expression: %s)", value, regionRegexp)
	}

	if parsedARN.AccountID != "" && !accountIDRegexp.MatchString(parsedARN.AccountID) {
		return fmt.Errorf("(%s) is an invalid ARN: invalid account ID value (expecting to match regular expression: %s)", value, accountIDRegexp)
	}

	if parsedARN.Resource == "" {
		return fmt.Errorf("(%s) is an invalid ARN: missing resource value", value)
	}

	return nil
}
