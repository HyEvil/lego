package utils

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/pkg/errors"
)

func GetPos(ctx antlr.ParserRuleContext) string {
	return fmt.Sprintf("at pos %d:%d", ctx.GetStart().GetLine(), ctx.GetStart().GetColumn())
}

func WithPosError(ctx antlr.ParserRuleContext, format string, args ...interface{}) error {
	return errors.WithMessagef(errors.Errorf("error at file %d:%d", ctx.GetStart().GetLine(), ctx.GetStart().GetColumn()), format, args...)
}
