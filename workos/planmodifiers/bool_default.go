package planmodifiers

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ planmodifier.Bool = &BoolDefaultModifier{}
)

type BoolDefaultModifier struct {
	Default bool
}

func (m BoolDefaultModifier) Description(_ context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to %v", m.Default)
}

func (m BoolDefaultModifier) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to `%v`", m.Default)
}

func (m BoolDefaultModifier) PlanModifyBool(_ context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if req.PlanValue.IsNull() {
		resp.PlanValue = types.BoolValue(m.Default)
	}
}

func BoolDefault(defaultValue bool) BoolDefaultModifier {
	return BoolDefaultModifier{
		Default: defaultValue,
	}
}
