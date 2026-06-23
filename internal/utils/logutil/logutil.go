package logutil

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/pkg/klog"
)

func ConfigureRuntime() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	hlog.SetOutput(os.Stdout)
	hlog.SetLevel(hlog.LevelWarn)

	klog.SetOutput(os.Stdout)
	klog.SetLevel(klog.LevelWarn)
}

func Duration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%.1fms", float64(d)/float64(time.Millisecond))
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}

func NormalizeUserAgent(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "-"
	}
	if len(value) > 64 {
		return value[:61] + "..."
	}
	return value
}

func NormalizeErr(err error) string {
	if err == nil {
		return "-"
	}
	return err.Error()
}
