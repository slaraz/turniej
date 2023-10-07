package turtles

import (
	"log"
	"reflect"
	"testing"
)

func Test_getCleanBoard(t *testing.T) {
	log.Println(getCleanBoard())
}

func Test_findPawn(t *testing.T) {
	type args struct {
		pawn  Color
		board []Field
	}
	tests := []struct {
		name            string
		args            args
		wantFieldNumber int
		wantPawnNumber  int
	}{
		{
			name: "find pawn",
			args: args{
				pawn: Red,
				board: []Field{
					{
						Pawns: []Color{Yellow, Green},
					},
					{
						Pawns: []Color{Red, Blue},
					},
				},
			},
			wantFieldNumber: 1,
			wantPawnNumber:  0,
		}, {
			name: "find pawn 1",
			args: args{
				pawn: Red,
				board: []Field{
					{
						Pawns: []Color{},
					},
					{
						Pawns: []Color{Yellow},
					},
					{
						Pawns: []Color{Blue},
					},
					{
						Pawns: []Color{Green, Red},
					},
				},
			},
			wantFieldNumber: 3,
			wantPawnNumber:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFieldNumber, gotPawnNumber := findPawn(tt.args.pawn, tt.args.board)
			if gotFieldNumber != tt.wantFieldNumber {
				t.Errorf("findPawn() gotFieldNumber = %v, want %v", gotFieldNumber, tt.wantFieldNumber)
			}
			if gotPawnNumber != tt.wantPawnNumber {
				t.Errorf("findPawn() gotPawnNumber = %v, want %v", gotPawnNumber, tt.wantPawnNumber)
			}
		})
	}
}

func TestMovePawn(t *testing.T) {
	type args struct {
		board []Field
		pawn  Color
		move  int
	}
	tests := []struct {
		name    string
		args    args
		want    []Field
		wantErr bool
	}{
		{
			name: "Start",
			args: args{
				board: []Field{
					{Pawns: []Color{}},
					{Pawns: []Color{}},
					{Pawns: []Color{}}},
				pawn: Green,
				move: 1,
			},
			want:    []Field{{Pawns: []Color{Green}}, {Pawns: []Color{}}, {Pawns: []Color{}}},
			wantErr: false,
		},
		{
			name: "Move pawn one foward",
			args: args{
				board: []Field{
					{Pawns: []Color{}},
					{Pawns: []Color{Red}},
					{Pawns: []Color{}}},
				pawn: Red,
				move: 1,
			},
			want:    []Field{{Pawns: []Color{}}, {Pawns: []Color{}}, {Pawns: []Color{Red}}},
			wantErr: false,
		},
		{
			name: "Move pawn two pawns together foward",
			args: args{
				board: []Field{
					{Pawns: []Color{}},
					{Pawns: []Color{Red, Green}},
					{Pawns: []Color{}}},
				pawn: Red,
				move: 1,
			},
			want:    []Field{{Pawns: []Color{}}, {Pawns: []Color{}}, {Pawns: []Color{Red, Green}}},
			wantErr: false,
		},
		{
			name: "Move one pawn of two  foward",
			args: args{
				board: []Field{
					{Pawns: []Color{}},
					{Pawns: []Color{Red, Green}},
					{Pawns: []Color{}}},
				pawn: Green,
				move: 1,
			},
			want:    []Field{{Pawns: []Color{}}, {Pawns: []Color{Red}}, {Pawns: []Color{Green}}},
			wantErr: false,
		},
		{
			name: "Move two pawns of three  foward",
			args: args{
				board: []Field{
					{Pawns: []Color{}},
					{Pawns: []Color{Red, Green, Blue}},
					{Pawns: []Color{}}},
				pawn: Green,
				move: 1,
			},
			want:    []Field{{Pawns: []Color{}}, {Pawns: []Color{Red}}, {Pawns: []Color{Green, Blue}}},
			wantErr: false,
		},
		{
			name: "Move pawn on top of other pawn",
			args: args{
				board: []Field{
					{Pawns: []Color{Red}},
					{Pawns: []Color{Green, Blue}},
					{Pawns: []Color{}}},
				pawn: Red,
				move: 1,
			},
			want:    []Field{{Pawns: []Color{}}, {Pawns: []Color{Green, Blue, Red}}, {Pawns: []Color{}}},
			wantErr: false,
		},
		{
			name: "Pawn out of board right side",
			args: args{
				board: []Field{
					{Pawns: []Color{Red}},
					{Pawns: []Color{Green, Blue}},
					{Pawns: []Color{}}},
				pawn: Red,
				move: 4,
			},
			want: []Field{
				{Pawns: []Color{}},
				{Pawns: []Color{Green, Blue}},
				{Pawns: []Color{Red}}},
			wantErr: false,
		},
		{
			name: "Pawn out of board left side",
			args: args{
				board: []Field{
					{Pawns: []Color{Red}},
					{Pawns: []Color{Green, Blue}},
					{Pawns: []Color{}}},
				pawn: Blue,
				move: -4,
			},
			want: []Field{
				{Pawns: []Color{Red}},
				{Pawns: []Color{Green}},
				{Pawns: []Color{}}},
			wantErr: false,
		},
		{
			name: "Pawn out of board left side and pawn is out of board",
			args: args{
				board: []Field{
					{Pawns: []Color{Red}},
					{Pawns: []Color{}},
					{Pawns: []Color{}}},
				pawn: Blue,
				move: -4,
			},
			want: []Field{
				{Pawns: []Color{Red}},
				{Pawns: []Color{}},
				{Pawns: []Color{}}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MovePawn(tt.args.board, tt.args.pawn, tt.args.move)
			if (err != nil) != tt.wantErr {
				t.Errorf("MovePawn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MovePawn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckIfGameOver(t *testing.T) {
	type args struct {
		board []Field
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 Color
	}{
		{
			name: "game isn't over",
			args: args{
				board: []Field{
					{Pawns: []Color{Red}},
					{Pawns: []Color{Green, Blue}},
					{Pawns: []Color{}}},
			},
			want:  false,
			want1: "",
		},
		{
			name: "game is over with one pawn",
			args: args{
				board: []Field{
					{Pawns: []Color{}},
					{Pawns: []Color{Blue}},
					{Pawns: []Color{Green}}},
			},
			want:  true,
			want1: Green,
		},
		{
			name: "game is over by the bottom pawn",
			args: args{
				board: []Field{
					{Pawns: []Color{Red}},
					{Pawns: []Color{}},
					{Pawns: []Color{Green, Blue}}},
			},
			want:  true,
			want1: Green,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := CheckIfGameOver(tt.args.board)
			if got != tt.want {
				t.Errorf("CheckIfGameOver() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("CheckIfGameOver() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
