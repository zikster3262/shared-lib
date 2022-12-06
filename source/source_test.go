package source

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestGetAllSources(t *testing.T) {

	mockDB, _, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sources")

	type args struct {
		db *sqlx.DB
	}
	tests := []struct {
		name    string
		args    args
		want    []SQL
		wantErr bool
	}{
		{name: "test get all query",
			args: args{
				db: sqlxDB,
			},
			want: []SQL{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAllSources(tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllSources() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllSources() = %v, want %v", got, tt.want)
			}
		})
	}
}
