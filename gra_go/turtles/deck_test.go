package turtles

import (
	"log"
	"reflect"
	"testing"
)

func TestShuffle(t *testing.T) {

	d := shuffleDefaultDeck()
	if reflect.DeepEqual(DefaultDeck, d) {
		t.Errorf("The same collection Shuffle(%v) = %v", DefaultDeck, d)
	}
	if len(DefaultDeck) != len(d) {
		t.Errorf("Diffrent number od cards Shuffle(%v) = %v", DefaultDeck, d)
	}
	disitnct := make(map[int]bool)
	for _, c := range d {
		if disitnct[c.id] {
			t.Errorf("Duplicate card %v", c)
		}
		disitnct[c.id] = true
	}
	log.Println(d)
}

func TestDeck_GetCardFromDeck(t *testing.T) {
	tests := []struct {
		name     string
		deck     Deck
		want     Card
		wantDeck Deck
		wantErr  bool
	}{
		{
			name:     "",
			deck:     []Card{{id: 1}, {id: 2}, {id: 3}},
			want:     Card{id: 1},
			wantDeck: []Card{{id: 2}, {id: 3}},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.deck.GetCardFromDeck()
			if (err != nil) != tt.wantErr {
				t.Errorf("Deck.GetCardFromDeck() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Deck.GetCardFromDeck() = %v, want %v", got, tt.want)
			}

			if !reflect.DeepEqual(tt.deck, tt.wantDeck) {
				t.Errorf("Deck.GetCardFromDeck() = %v, want %v", tt.deck, tt.wantDeck)
			}
		})
	}
}
