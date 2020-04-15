package logic_test

import (
	"github.com/david-sorm/goblog/article/logic"
	"github.com/david-sorm/goblog/article/store"
	"github.com/david-sorm/goblog/article/store/mock"
	"github.com/david-sorm/goblog/article/store/postgres"
	"reflect"
	"testing"
)

func TestParseArticleStore(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want store.ArticleStore
	}{
		{
			name: "mock store",
			args: args{str: "mock"},
			want: &mock.MockStore{},
		},
		{
			name: "postgres store",
			args: args{str: "postgres"},
			want: &postgres.PostgresStore{},
		},
		{
			name: "invalid store",
			args: args{str: "this store shouldn't exist"},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := logic.ParseArticleStore(tt.args.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseArticleStore() = %v, want %v", got, tt.want)
			}
		})
	}
}
