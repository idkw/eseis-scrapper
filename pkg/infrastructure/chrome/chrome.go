package chrome

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/sirupsen/logrus"
)

type Chrome struct {
	ctx   *context.Context
	close func()
}

func NewChrome() *Chrome {
	ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithLogf(logrus.Infof))
	return &Chrome{
		ctx: &ctx,
		close: func() {
			cancel()
		},
	}
}

func (c *Chrome) RunTasks(tasks chromedp.Tasks) error {
	if err := chromedp.Run(*c.ctx, tasks); err != nil {
		return err
	}
	return nil
}

func (c *Chrome) Close() {
	c.close()
}
