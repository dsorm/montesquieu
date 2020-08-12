package logic_test

import (
	"github.com/david-sorm/goblog/article/logic"
	"github.com/david-sorm/goblog/store"
	"github.com/david-sorm/goblog/store/mock"
	"github.com/david-sorm/goblog/store/postgres"
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
			name: "mock store",
			args: args{str: "mock"},
			want: &mock.Store{},
		},
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
