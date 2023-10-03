package turtles

import (
	"reflect"
	"testing"
)

func Test_checkIfExist(t *testing.T) {
	type args struct {
		in0 []Card
		in1 Card
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Card exist",
			args: args{
				in0: []Card{{Color: "red", Move: 1}, {Color: "blue", Move: 2}},
				in1: Card{Color: "red", Move: 1},
			},
			want: true,
		}, {
			name: "Card doesn't exist",
			args: args{
				in0: []Card{{Color: "red", Move: 1}, {Color: "blue", Move: 2}},
				in1: Card{Color: "green", Move: 1},
			},
			want: false,
		}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkIfExist(tt.args.in0, tt.args.in1); got != tt.want {
				t.Errorf("checkIfExist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_removeCard(t *testing.T) {
	type args struct {
		cards []Card
		c     Card
	}
	tests := []struct {
		name string
		args args
		want []Card
	}{
		{
			name: "card removed",
			args: args{
				cards: []Card{{Color: "red", Move: 1}, {Color: "blue", Move: 2}, {Color: "green", Move: 3}},
				c:     Card{Color: "blue", Move: 2},
			},
			want: []Card{{Color: "red", Move: 1}, {Color: "green", Move: 3}},
		},
		{
			name: "first removed",
			args: args{
				cards: []Card{{Color: "red", Move: 1}, {Color: "blue", Move: 2}, {Color: "green", Move: 3}},
				c:     Card{Color: "red", Move: 1},
			},
			want: []Card{{Color: "blue", Move: 2}, {Color: "green", Move: 3}},
		},
		{
			name: "last removed",
			args: args{
				cards: []Card{{Color: "red", Move: 1}, {Color: "blue", Move: 2}, {Color: "green", Move: 3}},
				c:     Card{Color: "green", Move: 3},
			},
			want: []Card{{Color: "red", Move: 1}, {Color: "blue", Move: 2}},
		},
		{
			name: "no removed",
			args: args{
				cards: []Card{{Color: "red", Move: 1}, {Color: "blue", Move: 2}, {Color: "green", Move: 3}},
				c:     Card{Color: "purpuble", Move: 1},
			},
			want: []Card{{Color: "red", Move: 1}, {Color: "blue", Move: 2}, {Color: "green", Move: 3}},
		},
		{
			name: "last one removed",
			args: args{
				cards: []Card{{Color: "red", Move: 1}},
				c:     Card{Color: "red", Move: 1},
			},
			want: []Card{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeCard(tt.args.cards, tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeCard() = %v, want %v", got, tt.want)
			}
		})
	}
}
