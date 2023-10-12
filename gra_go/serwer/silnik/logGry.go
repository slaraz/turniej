package silnik

import (
	"encoding/json"

	"github.com/slaraz/turniej/gra_go/proto"
)

type logGry struct {
	NowaGra        logNowaGra
	PrzebiegGry    []logStatus
	WynikGry       logWynikGry
	mapaKolorGracz map[proto.KolorZolwia]*gracz
	mapaNumerGracz map[int]*gracz
}

type logNowaGra struct {
	GraID  string
	Gracze []logGracz
}

type logGracz struct {
	NumerGracza int
	NazwaGracza string
	KolorGracza string
}

type logStatus struct {
	NazwaGracza  string
	TwojKolor    string
	TwojeKarty   []string
	Plansza      []logPole
	ZagraneKatry []logZagranaKarta
}

type logPole struct {
	Zolwie []string
}

type logZagranaKarta struct {
	NazwaGracza  string
	ZagranaKarta string
}

type logWynikGry struct {
	WygranyGracz string
}

func nowyLog(graID string) *logGry {
	return &logGry{
		NowaGra: logNowaGra{
			GraID:  graID,
			Gracze: []logGracz{},
		},
		PrzebiegGry:    []logStatus{},
		WynikGry:       logWynikGry{},
		mapaKolorGracz: map[proto.KolorZolwia]*gracz{},
		mapaNumerGracz: map[int]*gracz{},
	}
}

func (l *logGry) dodajGraczy(gracze []*gracz) {
	for i, gracz := range gracze {
		lg := logGracz{
			NumerGracza: i + 1,
			NazwaGracza: gracz.nazwaGracza,
			KolorGracza: proto.KolorZolwia_name[int32(gracz.kolorGracza)],
		}
		l.NowaGra.Gracze = append(l.NowaGra.Gracze, lg)
	}
	for _, gracz := range gracze {
		l.mapaKolorGracz[gracz.kolorGracza] = gracz
	}
	for _, gracz := range gracze {
		l.mapaNumerGracz[gracz.numerGracza] = gracz
	}
}

func (l *logGry) dodajStan(stan *proto.StanGry) {
	ls := logStatus{
		NazwaGracza:  l.mapaKolorGracz[stan.TwojKolor].nazwaGracza,
		TwojKolor:    stan.TwojKolor.String(),
		Plansza:      logPola(stan.Plansza),
		TwojeKarty:   logKarty(stan.TwojeKarty),
		ZagraneKatry: logZagraneKarty(stan.ZagraneKarty, l),
	}
	l.PrzebiegGry = append(l.PrzebiegGry, ls)
}

func logPola(plansza []*proto.Pole) []logPole {
	pola := []logPole{}
	for _, pole := range plansza {
		zolwie := []string{}
		for _, z := range pole.Zolwie {
			zolwtxt := proto.KolorZolwia_name[int32(z)]
			zolwie = append(zolwie, zolwtxt)
		}
		lp := logPole{
			Zolwie: zolwie,
		}
		pola = append(pola, lp)
	}
	return pola
}

func logKarty(karty []proto.Karta) []string {
	kartytxt := []string{}
	for _, karta := range karty {
		k := proto.Karta_name[int32(karta)]
		kartytxt = append(kartytxt, k)
	}
	return kartytxt
}

func logZagraneKarty(zagraneKarty []*proto.ZagranaKarta, l *logGry) []logZagranaKarta {
	zagrane := []logZagranaKarta{}
	for _, zg := range zagraneKarty {
		zagrana := logZagranaKarta{
			NazwaGracza:  l.mapaNumerGracz[int(zg.NumerGracza)].nazwaGracza,
			ZagranaKarta: zg.Karta.String(),
		}
		zagrane = append(zagrane, zagrana)
	}
	return zagrane
}

func (l *logGry) dodajKoniec(stan *proto.StanGry) {
	l.WynikGry.WygranyGracz = l.mapaNumerGracz[int(stan.KtoWygral)].nazwaGracza
}

func (l *logGry) getJSON() string {
	j, _ := json.MarshalIndent(l, "", "  ")
	return string(j)
}
