package logic_test

import (
	"github.com/david-sorm/montesquieu/article/logic"
	"github.com/david-sorm/montesquieu/store"
	"github.com/david-sorm/montesquieu/store/postgres"
	"reflect"
	"testing"
)

func TestParseStore(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want store.Store
	}{
		{
			name: "postgres store",
			args: args{str: "postgres"},
			want: &postgres.Store{},
		},
		{
			name: "invalid store",
			args: args{str: "this store shouldn't exist"},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := logic.ParseStore(tt.args.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseStore() = %v, want %v", got, tt.want)
			}
		})
	}
}
