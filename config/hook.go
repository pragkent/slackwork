package config

type Slack struct {
	Hooks []*Hook `yaml:"hooks"`
}

func (s *Slack) GetHook(name string) *Hook {
	for _, h := range s.Hooks {
		if h.Name == name {
			return h
		}
	}

	return nil
}

type Hook struct {
	Name    string `yaml:"name"`
	Secret  string `yaml:"secret"`
	AgentID int    `yaml:"agent_id"`
}

func (h Hook) Auth(secret string) bool {
	return h.Secret == secret
}
