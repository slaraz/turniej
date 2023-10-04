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
			want1: 2,
		},
		{
			name: "Red Wins",
			args: args{
				board:   []Field{{}, {}, {}, {}, {Pawns: []Color{"purpule", "red"}}},
				players: []Player{{Color: "red"}, {Color: "blue"}, {Color: "green"}},
			},
			want:  Player{Color: "red"},
			want1: 0,
		},
		{
			name: "Blue Wins",
			args: args{
				board:   []Field{{}, {}, {}, {Pawns: []Color{"blue"}}, {Pawns: []Color{"purpule", "red"}}},
				players: []Player{{Color: "yellow"}, {Color: "blue"}},
			},
			want:  Player{Color: "blue"},
			want1: 1,
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
