package fusefrontend

import (
	"github.com/rfjakob/gocryptfs/v2/internal/tlog"
	"os"
	"strings"
)

// Of converts variant T to its slice
func Of[T any](s ...T) []T {
	return s
}

// Tracef traces logging.
func Tracef(tags []string, format string, v ...interface{}) {
	tracer(tags).Printf(format, v...)
}

type PrintfAware interface {
	Printf(format string, v ...interface{})
}

type noopTracer struct{}

func (n noopTracer) Printf(format string, v ...interface{}) {}

var noopTrace = &noopTracer{}

var tracer = func() func([]string) PrintfAware {
	traceEnv := os.Getenv("TRACE")
	if traceEnv == "" {
		return func([]string) PrintfAware {
			return noopTrace
		}
	}

	level := traceEnv
	allowTags := ""
	if p := strings.IndexByte(traceEnv, ':'); p >= 0 {
		level = traceEnv[:p]
		allowTags = traceEnv[p+1:]
	}

	switch strings.ToLower(level) {
	case "warn":
		if allowTags == "" {
			return func(tags []string) PrintfAware {
				return tlog.Warn
			}
		} else {
			allowTagsSlice := strings.Split(allowTags, ",")
			return func(tags []string) PrintfAware {
				if containsAny(allowTagsSlice, tags) {
					return tlog.Warn
				}
				return noopTrace
			}
		}
	default: // info
		if allowTags == "" {
			return func(tags []string) PrintfAware {
				return tlog.Info
			}
		} else {
			allowTagsSlice := strings.Split(allowTags, ",")
			return func(tags []string) PrintfAware {
				if containsAny(allowTagsSlice, tags) {
					return tlog.Info
				}
				return noopTrace
			}
		}

	}
}()

func containsAny(src, subs []string) bool {
	for _, sr := range src {
		for _, su := range subs {
			if sr == su {
				return true
			}
		}
	}

	return false
}
