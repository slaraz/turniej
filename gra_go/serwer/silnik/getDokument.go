package silnik

import "time"

func (g *gra) getDokument() map[string]any {
	dokument := map[string]any{
		"timestamp": time.Now(),
		"graID":     g.graID,
		"logGry":    g.logGry,
		"gracze":    g.getGracze(),
	}
	return dokument
}

func (g *gra) getGracze() []string {
	//for g.gracze
	return []string{"gracz1", "gracz2"}
}
