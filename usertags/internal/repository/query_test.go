package repository

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrepareBulkInsertStmt(t *testing.T) {
	type args struct {
		rowsL     int
		tableName string
		cols      []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "success",
			args: args{
				rowsL:     3,
				tableName: "tag",
				cols:      []string{"keyword", "created_at", "updated_at"},
			},
			want: "INSERT INTO tag (keyword,created_at,updated_at) VALUES ($1,$2,$3),($4,$5,$6),($7,$8,$9)",
		}, {
			name: "success",
			args: args{
				rowsL:     0,
				tableName: "user_tags",
				cols:      []string{"user_id", "tag_id"},
			},
			want: "INSERT INTO user_tags (user_id,tag_id) ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := buildInsertStr(tt.args.rowsL, tt.args.tableName, tt.args.cols)
			assert.Equal(t, tt.want, got, "received statement doesn't match the expected one")

		})
	}
}
