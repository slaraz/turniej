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
				in0: []Card{{color: "red", move: 1}, {color: "blue", move: 2}},
				in1: Card{color: "red", move: 1},
			},
			want: true,
		}, {
			name: "Card doesn't exist",
			args: args{
				in0: []Card{{color: "red", move: 1}, {color: "blue", move: 2}},
				in1: Card{color: "green", move: 1},
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
				cards: []Card{{color: "red", move: 1}, {color: "blue", move: 2}, {color: "green", move: 3}},
				c:     Card{color: "blue", move: 2},
			},
			want: []Card{{color: "red", move: 1}, {color: "green", move: 3}},
		},
		{
			name: "first removed",
			args: args{
				cards: []Card{{color: "red", move: 1}, {color: "blue", move: 2}, {color: "green", move: 3}},
				c:     Card{color: "red", move: 1},
			},
			want: []Card{{color: "blue", move: 2}, {color: "green", move: 3}},
		},
		{
			name: "last removed",
			args: args{
				cards: []Card{{color: "red", move: 1}, {color: "blue", move: 2}, {color: "green", move: 3}},
				c:     Card{color: "green", move: 3},
			},
			want: []Card{{color: "red", move: 1}, {color: "blue", move: 2}},
		},
		{
			name: "no removed",
			args: args{
				cards: []Card{{color: "red", move: 1}, {color: "blue", move: 2}, {color: "green", move: 3}},
				c:     Card{color: "purpuble", move: 1},
			},
			want: []Card{{color: "red", move: 1}, {color: "blue", move: 2}, {color: "green", move: 3}},
		},
		{
			name: "last one removed",
			args: args{
				cards: []Card{{color: "red", move: 1}},
				c:     Card{color: "red", move: 1},
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
