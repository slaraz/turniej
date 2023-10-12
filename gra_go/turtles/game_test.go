package turtles

import (
	"reflect"
	"testing"
)

func Test_findWinner(t *testing.T) {
	type args struct {
		board   []Field
		players []Player
	}
	tests := []struct {
		name  string
		args  args
		want  Player
		want1 int
	}{
		{
			name: "Green Wins",
			args: args{
				board:   []Field{{}, {}, {}, {}, {Pawns: []Color{"green"}}},
				players: []Player{{Color: "red"}, {Color: "blue"}, {Color: "green"}},
			},
			want:  Player{Color: "green"},
			want1: 3,
		},
		{
			name: "Red Wins",
			args: args{
				board:   []Field{{}, {}, {}, {}, {Pawns: []Color{"purpule", "red"}}},
				players: []Player{{Color: "red"}, {Color: "blue"}, {Color: "green"}},
			},
			want:  Player{Color: "red"},
			want1: 1,
		},
		{
			name: "Blue Wins",
			args: args{
				board:   []Field{{}, {}, {}, {Pawns: []Color{"blue"}}, {Pawns: []Color{"purpule", "red"}}},
				players: []Player{{Color: "yellow"}, {Color: "blue"}},
			},
			want:  Player{Color: "blue"},
			want1: 2,
		},
		{
			name: "blue Wins",
			args: args{
				board:   []Field{{}, {}, {}, {Pawns: []Color{"red"}}, {Pawns: []Color{"yellow", "blue"}}},
				players: []Player{{Color: "yellow"}, {Color: "blue"}},
			},
			want:  Player{Color: "blue"},
			want1: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := findWinner(tt.args.board, tt.args.players)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findWinner() = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("findWinner() = %v, want %v", got1, tt.want1)
			}

		})
	}
}

func Test_shuffleColorsd(t *testing.T) {
	type args struct {
		colors []Color
	}
	tests := []struct {
		name string
		args args
		want []Color
	}{
		{
			name: "collors are shuffled",
			args: args{
				colors: []Color{"red", "green", "blue", "yellow"},
			},

			want: []Color{"red", "green", "blue", "yellow"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shuffleColors(tt.args.colors); reflect.DeepEqual(got, tt.want) {
				t.Errorf("not shuffling() = %v, want diffrent %v", got, tt.want)
			}
		})
	}
}
