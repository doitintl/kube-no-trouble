package printer

import (
	"context"
	"fmt"
	"strings"

	ctxKey "github.com/doitintl/kube-no-trouble/pkg/context"
)

func mapToCommaSeparatedString(m map[string]interface{}) string {
	var sb strings.Builder
	for k, v := range m {
		if sb.Len() > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%s:%v", k, v))
	}
	return sb.String()
}

func shouldShowLabels(ctx context.Context) (*bool, error) {
	if v := ctx.Value(ctxKey.LABELS_CTX_KEY); v != nil {
		return ctx.Value(ctxKey.LABELS_CTX_KEY).(*bool), nil
	}
	return nil, fmt.Errorf("labels flag not present in the context")
}
