package postgres

import (
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
)

func mockDB() (mock sqlmock.Sqlmock, err error) {
	db, mock, err = sqlmock.New()
	return
}

func Test_ListTables(t *testing.T) {
	mock, err := mockDB()
	if err != nil {
		t.Errorf("expected no errors, but got %v", err)
	}
	mock.ExpectQuery("SELECT").WillReturnRows(
		sqlmock.NewRows([]string{"tablename"}).AddRow("test"))

	table, err := ListTables()
	if err != nil {
		t.Errorf("expected no errors, but got %v", err)
	}
	if len(table) != 1 {
		t.Errorf("expected 1, but got %v", len(table))
	}
	if table[0].Name != "test" {
		t.Errorf("expected \"test\", but got %v", table[0].Name)
	}

	mock.ExpectQuery("SELECT").WillReturnError(errors.New("error test"))
	_, err = ListTables()
	if err == nil {
		t.Errorf("expected errors, but got nil")
	}

}

func Test_LoadTableJSON(t *testing.T) {
	mock, err := mockDB()
	if err != nil {
		t.Errorf("expected no errors, but got %v", err)
	}
	mock.ExpectQuery(`SELECT`).WillReturnRows(
		sqlmock.NewRows([]string{"field"}).AddRow("test"))

	j, err := LoadTableJSON("tablename")
	if err != nil {
		t.Errorf("expected no errors, but got %v", err)
	}
	expected := "[\n\t{\n\t\t\"field\": \"test\"\n\t}\n]"
	if string(j) != expected {
		t.Errorf("expected %q, but got %q", expected, string(j))
	}

	mock.ExpectQuery("SELECT").WillReturnError(errors.New("error test"))
	_, err = LoadTableJSON("tablename")
	if err == nil {
		t.Errorf("expected errors, but got nil")
	}
}

func Test_LoadTableCSV(t *testing.T) {
	mock, err := mockDB()
	if err != nil {
		t.Errorf("expected no errors, but got %v", err)
	}
	mock.ExpectQuery(`SELECT`).WillReturnRows(
		sqlmock.NewRows([]string{"field"}).AddRow("test"))

	j, err := LoadTableCSV("tablename")
	if err != nil {
		t.Errorf("expected no errors, but got %v", err)
	}
	expected := "field\ntest\n"
	if string(j) != expected {
		t.Errorf("expected %q, but got %q", expected, string(j))
	}

	mock.ExpectQuery("SELECT").WillReturnError(errors.New("error test"))
	_, err = LoadTableCSV("tablename")
	if err == nil {
		t.Errorf("expected errors, but got nil")
	}
}

func Test_scanner_Scan(t *testing.T) {
	type fields struct {
		value interface{}
	}
	type args struct {
		src interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "int64",
			fields: fields{
				value: 10,
			},
			args: args{
				src: 10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := &scanner{
				value: tt.fields.value,
			}
			if err := scanner.Scan(tt.args.src); (err != nil) != tt.wantErr {
				t.Errorf("scanner.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
