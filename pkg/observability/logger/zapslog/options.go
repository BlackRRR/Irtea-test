package zapslog

import "log/slog"

// A HandlerOption configures a slog Handler.
type HandlerOption interface {
	apply(*Handler)
}

// handlerOptionFunc wraps a func so it satisfies the Option interface.
type handlerOptionFunc func(*Handler)

func (f handlerOptionFunc) apply(handler *Handler) {
	f(handler)
}

// WithName configures the Logger to annotate each message with the logger name.
func WithName(name string) HandlerOption {
	return handlerOptionFunc(func(h *Handler) {
		h.name = name
	})
}

// WithCaller configures the Logger to include the filename and line number
// of the caller in log messages--if available.
func WithCaller(enabled bool) HandlerOption {
	return handlerOptionFunc(func(handler *Handler) {
		handler.addCaller = enabled
	})
}

// WithCallerSkip increases the number of callers skipped by caller annotation
// (as enabled by the [WithCaller] option).
//
// When building wrappers around the Logger,
// supplying this Option prevents Zap from always reporting
// the wrapper code as the caller.
func WithCallerSkip(skip int) HandlerOption {
	return handlerOptionFunc(func(log *Handler) {
		log.callerSkip += skip
	})
}

// AddStacktraceAt configures the Logger to record a stack trace
// for all messages at or above a given level.
func AddStacktraceAt(lvl slog.Level) HandlerOption {
	return handlerOptionFunc(func(log *Handler) {
		log.addStackAt = lvl
	})
}
