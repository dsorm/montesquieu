package config

import "testing"

func Test_file_configEmpty(t *testing.T) {
	type fields struct {
		BlogName         string
		ArticlesPerPage  string
		ListenOn         string
		Store            string
		StoreHost        string
		StoreDB          string
		StoreUser        string
		StorePassword    string
		CachingStore     string
		HotSwapTemplates string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "empty",
			fields: fields{
				BlogName:         "",
				ArticlesPerPage:  "",
				ListenOn:         "",
				Store:            "",
				StoreHost:        "",
				StoreDB:          "",
				StoreUser:        "",
				StorePassword:    "",
				CachingStore:     "",
				HotSwapTemplates: "",
			},
			want: true,
		},
		{
			name: "all full",
			fields: fields{
				BlogName:         "dsf",
				ArticlesPerPage:  "dfs",
				ListenOn:         "sds",
				Store:            "sdadw",
				StoreHost:        "sddsrw",
				StoreDB:          "sdsfs",
				StoreUser:        "sdsa",
				StorePassword:    "dasds",
				CachingStore:     "dasd",
				HotSwapTemplates: "sdfa",
			},
			want: false,
		},
		{
			name: "some full",
			fields: fields{
				BlogName:         "dsdasdas",
				ArticlesPerPage:  "",
				ListenOn:         "",
				Store:            "",
				StoreHost:        "dasdasds",
				StoreDB:          "",
				StoreUser:        "",
				StorePassword:    "dasddasd",
				CachingStore:     "",
				HotSwapTemplates: "",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &file{
				BlogName:         tt.fields.BlogName,
				ArticlesPerPage:  tt.fields.ArticlesPerPage,
				ListenOn:         tt.fields.ListenOn,
				Store:            tt.fields.Store,
				StoreHost:        tt.fields.StoreHost,
				StoreDB:          tt.fields.StoreDB,
				StoreUser:        tt.fields.StoreUser,
				StorePassword:    tt.fields.StorePassword,
				CachingStore:     tt.fields.CachingStore,
				HotSwapTemplates: tt.fields.HotSwapTemplates,
			}
			if got := cfg.configEmpty(); got != tt.want {
				t.Errorf("configEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}
