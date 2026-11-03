package zsms

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gmode"
)

type mockDriver struct{}

// Send 本地联调驱动：不真正发短信，仅打印日志。
func (d *mockDriver) Send(ctx context.Context, msg Message) error {
	if gmode.IsDevelop() {
		g.Log().Infof(ctx, "zsms(mock): mobile=%s template=%s code=%s", maskMobile(msg.Mobile), msg.Template, msg.Code)
		return nil
	}
	g.Log().Infof(ctx, "zsms(mock): 短信未真实下发 mobile=%s template=%s", maskMobile(msg.Mobile), msg.Template)
	return fmt.Errorf("zsms: mock 驱动仅允许在 develop 模式使用，短信未真实下发")
}
