package config

const DefaultWorkDir = "/var/www"

type Config struct {
	Repository string `json:"repository"`
	WorkDir    string `json:"work_dir,omitempty"`
}

func (c Config) WorkDirectory() string {
	if c.WorkDir != "" {
		return c.WorkDir
	}
	return DefaultWorkDir
}
