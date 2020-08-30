package mock

import (
	"github.com/david-sorm/goblog/article"
	"github.com/david-sorm/goblog/store"
	"reflect"
	"testing"
)

// struct fields
type fields struct {
	cfg                 store.StoreConfig
	articlesByTimestamp []article.Article
	articlesByID        map[string]article.Article
}

// makes test struct field content for unit testing
func getTestFields() fields {
	f := fields{
		cfg: store.StoreConfig{
			ArticlesPerIndexPage: 3,
		},
		articlesByTimestamp: []article.Article{
			{
				Timestamp: 10,
				ID:        1,
				Title:     "Article 1",
				Content:   "This is Article 1.",
			},
			{
				Timestamp: 9,
				ID:        2,
				Title:     "Article 2",
				Content:   "This is Article 2.",
			},
			{
				Timestamp: 8,
				ID:        3,
				Title:     "Article 3",
				Content:   "This is Article 3.",
			},
			{
				Timestamp: 7,
				ID:        4,
				Title:     "Article 4",
				Content:   "This is Article 4.",
			},
			{
				Timestamp: 6,
				ID:        5,
				Title:     "Article 5",
				Content:   "This is Article 5.",
			},
			{
				Timestamp: 5,
				ID:        6,
				Title:     "Article 6",
				Content:   "This is Article 6.",
			},
			{
				Timestamp: 4,
				ID:        7,
				Title:     "Article 7",
				Content:   "This is Article 7.",
			},
			{
				Timestamp: 3,
				ID:        8,
				Title:     "Article 8",
				Content:   "This is Article 8.",
			},
			{
				Timestamp: 2,
				ID:        9,
				Title:     "Article 9",
				Content:   "This is Article 9.",
			},
			{
				Timestamp: 1,
				ID:        10,
				Title:     "Article 10",
				Content:   "This is Article 10.",
			},
		},
		articlesByID: map[string]article.Article{
			"1": {
				Timestamp: 10,
				ID:        1,
				Title:     "Article 1",
				Content:   "This is Article 1.",
			},
			"2": {
				Timestamp: 9,
				ID:        2,
				Title:     "Article 2",
				Content:   "This is Article 2.",
			},
			"3": {
				Timestamp: 8,
				ID:        3,
				Title:     "Article 3",
				Content:   "This is Article 3.",
			},
			"4": {
				Timestamp: 7,
				ID:        4,
				Title:     "Article 4",
				Content:   "This is Article 4.",
			},
			"5": {
				Timestamp: 6,
				ID:        5,
				Title:     "Article 5",
				Content:   "This is Article 5.",
			},
			"6": {
				Timestamp: 5,
				ID:        6,
				Title:     "Article 6",
				Content:   "This is Article 6.",
			},
			"7": {
				Timestamp: 4,
				ID:        7,
				Title:     "Article 7",
				Content:   "This is Article 7.",
			},
			"8": {
				Timestamp: 3,
				ID:        8,
				Title:     "Article 8",
				Content:   "This is Article 8.",
			},
			"9": {
				Timestamp: 2,
				ID:        9,
				Title:     "Article 9",
				Content:   "This is Article 9.",
			},
			"10": {
				Timestamp: 1,
				ID:        10,
				Title:     "Article 10",
				Content:   "This is Article 10.",
			},
		},
	}

	return f
}

func TestMockStore_GetArticleByID(t *testing.T) {

	type args struct {
		ID uint64
	}

	testFields := getTestFields()

	tests := []struct {
		name   string
		fields fields
		args   args
		want   article.Article
		want1  bool
	}{
		{
			name:   "Invalid Article",
			fields: testFields,
			args:   args{ID: 250604},
			want:   article.Article{},
			want1:  false,
		},
		{
			name:   "Article 1",
			fields: testFields,
			args:   args{ID: 1},
			want: article.Article{
				Timestamp: 10,
				ID:        1,
				Title:     "Article 1",
				Content:   "This is Article 1.",
			},
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &Store{
				cfg:                 tt.fields.cfg,
				articlesByTimestamp: tt.fields.articlesByTimestamp,
				articlesByID:        tt.fields.articlesByID,
			}
			got, got1 := ms.GetArticleByID(tt.args.ID)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetArticleByID() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetArticleByID() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMockStore_GetArticleNumber(t *testing.T) {
	testFields := getTestFields()
	tests := []struct {
		name   string
		fields fields
		want   uint64
	}{
		{
			name:   "getArticleNumber test",
			fields: testFields,
			want:   10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &Store{
				cfg:                 tt.fields.cfg,
				articlesByTimestamp: tt.fields.articlesByTimestamp,
				articlesByID:        tt.fields.articlesByID,
			}
			if got := ms.GetArticleNumber(); got != tt.want {
				t.Errorf("GetArticleNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMockStore_Init(t *testing.T) {
	type args struct {
		f   func()
		cfg store.StoreConfig
	}
	testFields := getTestFields()
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "init",
			fields: testFields,
			args: args{
				f:   func() {},
				cfg: testFields.cfg,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &Store{
				cfg:                 tt.fields.cfg,
				articlesByTimestamp: tt.fields.articlesByTimestamp,
				articlesByID:        tt.fields.articlesByID,
			}
			if err := ms.Init(tt.args.f, tt.args.cfg); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMockStore_LoadArticlesForIndex(t *testing.T) {
	type args struct {
		from uint64
		to   uint64
	}

	testFields := getTestFields()
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []article.Article
	}{
		{
			name:   "page 0 (first page) test",
			fields: testFields,
			args:   args{0, 3},
			want: []article.Article{
				{
					Timestamp: 10,
					ID:        1,
					Title:     "Article 1",
					Content:   "This is Article 1.",
				},
				{
					Timestamp: 9,
					ID:        2,
					Title:     "Article 2",
					Content:   "This is Article 2.",
				},
				{
					Timestamp: 8,
					ID:        3,
					Title:     "Article 3",
					Content:   "This is Article 3.",
				},
			},
		},
		{
			name:   "page 1 test",
			fields: testFields,
			args:   args{3, 6},
			want: []article.Article{
				{
					Timestamp: 7,
					ID:        4,
					Title:     "Article 4",
					Content:   "This is Article 4.",
				},
				{
					Timestamp: 6,
					ID:        5,
					Title:     "Article 5",
					Content:   "This is Article 5.",
				},
				{
					Timestamp: 5,
					ID:        6,
					Title:     "Article 6",
					Content:   "This is Article 6.",
				},
			},
		},
		{
			name:   "page 3 (last page) test",
			fields: testFields,
			args:   args{9, 10},
			want: []article.Article{
				{
					Timestamp: 1,
					ID:        10,
					Title:     "Article 10",
					Content:   "This is Article 10.",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &Store{
				cfg:                 tt.fields.cfg,
				articlesByTimestamp: tt.fields.articlesByTimestamp,
				articlesByID:        tt.fields.articlesByID,
			}
			if got := ms.LoadArticlesSortedByLatest(tt.args.from, tt.args.to); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadArticlesForIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}
