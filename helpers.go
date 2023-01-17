package main

import swissknife "github.com/Sagleft/swiss-knife"

func (b *bot) parseConfig() error {
	return swissknife.ParseStructFromJSONFile(configPath, &b.Config)
}
