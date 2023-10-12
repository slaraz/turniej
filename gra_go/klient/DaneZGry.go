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
	sg.GetTwojKolor()
}

func znajdzOstatnieZolwie(pole []*proto.Pole) []proto.KolorZolwia {
	ostatnieZolwie := []proto.KolorZolwia{
		proto.KolorZolwia_XXX,
		proto.KolorZolwia_RED, 
		proto.KolorZolwia_GREEN, 
		proto.KolorZolwia_BLUE,  
		proto.KolorZolwia_YELLOW,
		proto.KolorZolwia_PURPLE,
	}

	for _, p := range pole {
		if p != nil {
			ostatnieZolwie = p.Zolwie
			break
		}
	}

	return ostatnieZolwie
}

func naszePole(naszKolor proto.KolorZolwia, pole []*proto.Pole) int {
	for i, p := range pole {
		for _, z := range p.GetZolwie() {
			if z == naszKolor {
				return i
			}
		}
	}

	return 0
}
