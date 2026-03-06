package run

import (
	"errors"
	"fmt"
	"io"

	"github.com/expr-lang/expr"
	"github.com/suzuki-shunsuke/go-error-with-exit-code/ecerror"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

func testResult(stderr io.Writer, testCode string, result *TemplateInput) error {
	prog, err := expr.Compile(testCode, expr.Env(result), expr.AsBool())
	if err != nil {
		fmt.Fprintf(stderr, `[ERROR] compile an expression
%v`, err)
		return ecerror.Wrap(nil, 1)
	}
	output, err := expr.Run(prog, result)
	if err != nil {
		fmt.Fprintf(stderr, `[ERROR] evaluate an expression
%v`, err)
		return ecerror.Wrap(nil, 1)
	}
	f, ok := output.(bool)
	if !ok {
		return errors.New("the test result must be boolean")
	}
	if !f {
		return slogerr.With(errors.New("test failed"), "test", testCode) //nolint:wrapcheck
	}
	return nil
}
