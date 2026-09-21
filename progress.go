package dem

import "context"

func ReportProgress(fn ProgressFunc, stage string, done, total int) {
	if fn != nil {
		fn(stage, done, total)
	}
}

func CheckCtx(ctx context.Context) error {
	if ctx == nil {
		return nil
	}
	return ctx.Err()
}
