package main

import (
	"reflect"
	"testing"

	"github.com/slaraz/turniej/gra_go/proto"
)

func TestZnajdzOstatnieZolwie(t *testing.T) {
	type args struct {
		pole []*proto.Pole
	}
	tests := []struct {
		name string
		args args
		want []proto.KolorZolwia
	}{
		{
			name: "wszystkie pola puste",
			args: args{
				pole: make([]*proto.Pole, 10),
			},
			want: []proto.KolorZolwia{
				proto.KolorZolwia_XXX,
				proto.KolorZolwia_RED,
				proto.KolorZolwia_GREEN,
				proto.KolorZolwia_BLUE,
				proto.KolorZolwia_YELLOW,
				proto.KolorZolwia_PURPLE,
			},
		},
		{
			name: "zielony jedyny",
			args: args{
				pole: []*proto.Pole{
					{
						Zolwie: []proto.KolorZolwia{
							proto.KolorZolwia_GREEN,
						},
					},
				},
			},
			want: []proto.KolorZolwia{
				proto.KolorZolwia_GREEN,
			},
		},
		{
			name: "puste pole, zolty i czerwony ostatni, pozniej zielony",
			args: args{
				pole: []*proto.Pole{
					{
						Zolwie: []proto.KolorZolwia{},
					},
					{
						Zolwie: []proto.KolorZolwia{
							proto.KolorZolwia_YELLOW,
							proto.KolorZolwia_RED,
						},
					},
					{
						Zolwie: []proto.KolorZolwia{
							proto.KolorZolwia_GREEN,
						},
					},
				},
			},
			want: []proto.KolorZolwia{
				proto.KolorZolwia_YELLOW,
				proto.KolorZolwia_RED,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := znajdzOstatnieZolwie(tt.args.pole); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("znajdzOstatnieZolwie() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNaszePole(t *testing.T) {
	type args struct {
		naszKolor proto.KolorZolwia
		plansza   []*proto.Pole
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "pusta plansza",
			args: args{
				naszKolor: proto.KolorZolwia_GREEN,
				plansza: make([]*proto.Pole, 10),
			},
			want: -1,
		},
		{
			name: "plansza bez naszego koloru",
			args: args{
				naszKolor: proto.KolorZolwia_GREEN,
				plansza: []*proto.Pole{
					{
						Zolwie: []proto.KolorZolwia{
							proto.KolorZolwia_BLUE,
						},
					},
					{
						Zolwie: []proto.KolorZolwia{},
					},
					{
						Zolwie: []proto.KolorZolwia{
							proto.KolorZolwia_PURPLE,
							proto.KolorZolwia_YELLOW,
						},
					},
				},
			},
			want: -1,
		},
		{
			name: "na polu indeks 2",
			args: args{
				naszKolor: proto.KolorZolwia_GREEN,
				plansza: []*proto.Pole{
					{
						Zolwie: []proto.KolorZolwia{
							proto.KolorZolwia_BLUE,
						},
					},
					{
						Zolwie: []proto.KolorZolwia{
							proto.KolorZolwia_PURPLE,
						},
					},
					{
						Zolwie: []proto.KolorZolwia{
							proto.KolorZolwia_GREEN,
						},
					},
				},
			},
			want: 2,
		},
		{
			name: "na polu indeks 3, zestackowany z innym",
			args: args{
				naszKolor: proto.KolorZolwia_GREEN,
				plansza: []*proto.Pole{
					{
						Zolwie: []proto.KolorZolwia{},
					},
					{
						Zolwie: []proto.KolorZolwia{},
					},
					{
						Zolwie: []proto.KolorZolwia{},
					},
					{
						Zolwie: []proto.KolorZolwia{
							proto.KolorZolwia_BLUE,
							proto.KolorZolwia_GREEN,
							proto.KolorZolwia_PURPLE,
						},
					},
				},
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := naszePole(tt.args.naszKolor, tt.args.plansza); got != tt.want {
				t.Errorf("naszePole() = %v, want %v", got, tt.want)
			}
		})
	}
}
