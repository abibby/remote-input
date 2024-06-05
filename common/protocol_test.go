package common

import (
	"reflect"
	"testing"
)

func TestKeyEvent_UnmarshalBinary(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		k       *InputEvent
		args    args
		wantErr bool
		wantK   *InputEvent
	}{
		{
			name: "a",
			k:    &InputEvent{},
			args: args{[]byte{
				0x00,       // version
				0x00, 0x02, // key code
				0x00, 0x00, 0x00, 0x03, // flags
			}},
			wantK: &InputEvent{Key: 2, Flags: 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.k.UnmarshalBinary(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("KeyEvent.UnmarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.k, tt.wantK) {
				t.Errorf("KeyEvent.UnmarshalBinary() = %v, want %v", tt.k, tt.wantK)
			}
		})
	}
}

func TestKeyEvent_MarshalBinary(t *testing.T) {
	tests := []struct {
		name     string
		k        *InputEvent
		wantData []byte
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, err := tt.k.MarshalBinary()
			if (err != nil) != tt.wantErr {
				t.Errorf("KeyEvent.MarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("KeyEvent.MarshalBinary() = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}
