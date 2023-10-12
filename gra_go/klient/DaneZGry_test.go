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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := znajdzOstatnieZolwie(tt.args.pole); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("znajdzOstatnieZolwie() = %v, want %v", got, tt.want)
			}
		})
	}
}
