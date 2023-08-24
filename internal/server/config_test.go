package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestConfig_Load(t *testing.T) {
	type fields struct {
		path   string
		Logger *zap.Logger
		env    map[string]string
	}
	type args struct {
		envPath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "Load env file",
			fields: fields{
				Logger: zap.NewExample(),
			},
			args: args{
				envPath: "../../config/test.env",
			},
			want: map[string]string{
				//Database Info
				"DB_HOST":     "localhost",
				"DB_PORT":     "5555",
				"DB_USER":     "postgres",
				"DB_PASSWORD": "postgres",
				"DB_NAME":     "dvdrental",
				// HTTP Info
				"HTTP_PORT": "8080",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				path:   tt.fields.path,
				Logger: tt.fields.Logger,
				env:    tt.fields.env,
			}
			//TODO: supress/vallidate the log output
			got, err := c.Load(tt.args.envPath)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
