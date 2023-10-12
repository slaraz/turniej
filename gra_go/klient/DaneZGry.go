package main

import "github.com/slaraz/turniej/gra_go/proto"

type DaneZGry struct {
	OstatnieZolwie       []proto.KolorZolwia
	ZolwieNadNami        []proto.KolorZolwia
	ZolwiePodNami        []proto.KolorZolwia
	DomniemanyPrzeciwnik proto.KolorZolwia // gdy 1v1
	NaszePole            int
	KrokowDoKonca        int
}

func (dzg *DaneZGry) PobierzDaneZeStanuGry(sg *proto.StanGry) {
	dzg.NaszePole = naszePole(sg.TwojKolor, sg.Plansza)
	dzg.KrokowDoKonca = len(sg.Plansza) - dzg.NaszePole 

	dzg.OstatnieZolwie = znajdzOstatnieZolwie(sg.Plansza)
	dzg.ZolwiePodNami = znajdzZolwiePodNami(dzg.NaszePole, sg.TwojKolor, sg.Plansza)
	dzg.ZolwieNadNami = znajdzZolwieNadNami(dzg.NaszePole, sg.TwojKolor, sg.Plansza)
}

func znajdzZolwieNadNami(naszePole int, naszKolor proto.KolorZolwia, plansza []*proto.Pole) []proto.KolorZolwia {
	indeksNaszegoZolwia := 0

	for i, kolorZolwia := range plansza[naszePole].GetZolwie() {
		if kolorZolwia == naszKolor {
			indeksNaszegoZolwia = i
			break
		}
	}

	return plansza[naszePole].GetZolwie()[indeksNaszegoZolwia:]
}

func znajdzZolwiePodNami(naszePole int, naszKolor proto.KolorZolwia, plansza []*proto.Pole) []proto.KolorZolwia {
	indeksNaszegoZolwia := 0

	for i, kolorZolwia := range plansza[naszePole].GetZolwie() {
		if kolorZolwia == naszKolor {
			indeksNaszegoZolwia = i
			break
		}
	}

	return plansza[naszePole].GetZolwie()[:indeksNaszegoZolwia]
}

func znajdzOstatnieZolwie(plansza []*proto.Pole) []proto.KolorZolwia {
	ostatnieZolwie := []proto.KolorZolwia{
		proto.KolorZolwia_XXX,
		proto.KolorZolwia_RED, 
		proto.KolorZolwia_GREEN, 
		proto.KolorZolwia_BLUE,  
		proto.KolorZolwia_YELLOW,
		proto.KolorZolwia_PURPLE,
	}

	for _, p := range plansza {
		if len(p.GetZolwie()) != 0 {
			ostatnieZolwie = p.Zolwie
			break
		}
	}

	return ostatnieZolwie
}

func naszePole(naszKolor proto.KolorZolwia, plansza []*proto.Pole) int {
	for i, p := range plansza {
		for _, z := range p.GetZolwie() {
			if z == naszKolor {
				return i
			}
		}
	}

	return -1
}
