package sink

import (
	"fmt"

	"github.com/pragkent/slackwork/config"
)

type Controller struct {
	c     *config.Config
	sinks map[string]Sink
}

func NewController(c *config.Config) (*Controller, error) {
	ctrl := &Controller{
		c:     c,
		sinks: make(map[string]Sink),
	}

	if err := ctrl.initSinks(); err != nil {
		return nil, err
	}

	return ctrl, nil
}

func (c *Controller) initSinks() error {
	for _, hc := range c.c.Slack.Hooks {
		wc := c.c.WeWork.GetAgent(hc.AgentID)
		if wc == nil {
			return fmt.Errorf("Agent not found: %v", hc.AgentID)
		}

		c.sinks[hc.Name] = NewWeWork(c.c.WeWork.CorpID, wc.Secret, wc.ID)
	}

	return nil
}

func (c *Controller) Dispatch(name string, payload *Payload) error {
	sk, ok := c.sinks[name]
	if !ok {
		return fmt.Errorf("Sink not found: %v", name)
	}

	if err := sk.Send(payload); err != nil {
		return fmt.Errorf("Sink error: %v %v", name, err)
	}

	return nil
}
