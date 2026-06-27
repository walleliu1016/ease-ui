package app

import "github.com/akke/ease-ui/internal/settings"

func (a *App) GetSettings() (*settings.Config, error) {
	return settings.Load()
}

func (a *App) UpdateSettings(cfg *settings.Config) error {
	if err := settings.Save(cfg); err != nil {
		return err
	}
	// re-load + sync handler
	a.appMu.Lock()
	a.settings = cfg
	a.handler.SetAutoAllowBash(cfg.AutoAllowBash)
	a.appMu.Unlock()
	return nil
}
