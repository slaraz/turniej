package main

// func Test_getCardFromText(t *testing.T) {
// 	type args struct {
// 		text string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    game.Card
// 		wantErr bool
// 	}{
// 		{
// 			name: "methot work",
// 			args: args{
// 				text: "red,1",
// 			},
// 			want:    game.Card{Color: "red", Move: 1},
// 			wantErr: false,
// 		},
// 		{
// 			name: "methot work",
// 			args: args{
// 				text: " blue\r,2\r\n",
// 			},
// 			want:    game.Card{Color: "blue", Move: 2},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := getCardFromText(tt.args.text)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("getCardFromText() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("getCardFromText() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
