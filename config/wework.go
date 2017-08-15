package config

type WeWork struct {
	CorpID string   `yaml:"corp_id"`
	Agents []*Agent `yaml:"agents"`
}

func (w *WeWork) GetAgent(id int) *Agent {
	for _, a := range w.Agents {
		if a.ID == id {
			return a
		}
	}

	return nil
}

type Agent struct {
	ID     int    `yaml:"id"`
	Name   string `yaml:"name"`
	Secret string `yaml:"secret"`
}
