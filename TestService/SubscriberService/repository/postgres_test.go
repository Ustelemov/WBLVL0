package repository

import (
	"database/sql"
	"encoding/json"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ustelemov/WBLVL0/TestService/schema"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
)

func TestRepository_GetAllOrders(t *testing.T) {
	db, mock, err := sqlxmock.Newx()

	if err != nil {
		t.Fatalf("error while openning mock database connection: %s", err)
	}
	defer db.Close()

	s, err := &PostgresDB{DB: db, listener: nil, config: Config{}}, nil

	if err != nil {
		t.Fatalf("error while creating repository: %s", err)
	}

	res := []schema.OrderJSON{
		{
			JSON:     []byte(""),
			OrderUID: "91bf621d-aaf3-4a0e-9048-329720a49b50",
		},
	}

	testTable := []struct {
		name    string
		s       *PostgresDB
		mock    func()
		want    []schema.OrderJSON
		wantErr bool
	}{
		{
			name: "OK",
			s:    s,
			mock: func() {
				rows := sqlxmock.NewRows([]string{"uuid", "order_data"}).AddRow("91bf621d-aaf3-4a0e-9048-329720a49b50", []byte(""))
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM ordersjson")).WillReturnRows(rows)
			},
			want: res,
		},
		{
			name: "NoRows",
			s:    s,
			mock: func() {
				rows := sqlxmock.NewRows([]string{"uuid", "order_data"})
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM ordersjson")).WillReturnRows(rows)
			},
			want: make([]schema.OrderJSON, 0),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()
			arr, err := testCase.s.GetAllOrders()

			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, arr, testCase.want)
			}

		})
	}

}

func TestRepository_CheckExists(t *testing.T) {
	db, mock, err := sqlxmock.Newx()

	if err != nil {
		t.Fatalf("error while openning mock database connection: %s", err)
	}
	defer db.Close()

	s, err := &PostgresDB{DB: db, listener: nil, config: Config{}}, nil

	if err != nil {
		t.Fatalf("error while creating repository: %s", err)
	}

	testTable := []struct {
		name       string
		s          *PostgresDB
		inputOrder *schema.Order
		mock       func(*schema.Order)
		want       bool
		wantErr    bool
	}{
		{
			name:       "Exists",
			s:          s,
			inputOrder: &schema.Order{OrderUID: "91bf621d-aaf3-4a0e-9048-329720a49b50"},
			mock: func(ord *schema.Order) {
				rows := sqlxmock.NewRows([]string{"uuid", "order_data"}).AddRow("91bf621d-aaf3-4a0e-9048-329720a49b50", []byte(""))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM ordersjson WHERE uuid = $1`)).WithArgs(ord.OrderUID).WillReturnRows(rows)
			},
			want: true,
		},
		{
			name:       "NoExists",
			s:          s,
			inputOrder: &schema.Order{OrderUID: "91bf621d-aaf3-4a0e-9048-329720a49b50"},
			mock: func(ord *schema.Order) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM ordersjson WHERE uuid = $1`)).WithArgs(ord.OrderUID).WillReturnError(sql.ErrNoRows)
			},
			want: false,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock(testCase.inputOrder)
			res, err := testCase.s.CheckExists(testCase.inputOrder)

			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, res, testCase.want)
			}

		})
	}
}

func TestRepository_SaveOrderInRepository(t *testing.T) {
	db, mock, err := sqlxmock.Newx()

	if err != nil {
		t.Fatalf("error while openning mock database connection: %s", err)
	}
	defer db.Close()

	s, err := &PostgresDB{DB: db, listener: nil, config: Config{}}, nil

	if err != nil {
		t.Fatalf("error while creating repository: %s", err)
	}
	testTable := []struct {
		name       string
		s          *PostgresDB
		inputOrder *schema.Order
		mock       func(*schema.Order, []byte)
		wantErr    bool
	}{
		{
			name:       "OK",
			s:          s,
			inputOrder: &schema.Order{OrderUID: "91bf621d-aaf3-4a0e-9048-329720a49b50"},
			mock: func(ord *schema.Order, js []byte) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM ordersjson WHERE uuid = $1`)).WithArgs(ord.OrderUID).WillReturnError(sql.ErrNoRows)
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO ordersjson (uuid, order_data) VALUES ($1,$2)")).WithArgs(ord.OrderUID, js).WillReturnResult(sqlxmock.NewResult(1, 1))
			},
		},
		{
			name:       "RowExists",
			s:          s,
			inputOrder: &schema.Order{OrderUID: "91bf621d-aaf3-4a0e-9048-329720a49b50"},
			mock: func(ord *schema.Order, js []byte) {
				rows := sqlxmock.NewRows([]string{"uuid", "order_data"}).AddRow("91bf621d-aaf3-4a0e-9048-329720a49b50", []byte(""))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM ordersjson WHERE uuid = $1`)).WithArgs(ord.OrderUID).WillReturnRows(rows)
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO ordersjson (uuid, order_data) VALUES ($1,$2)")).WithArgs(ord.OrderUID, js).WillReturnResult(sqlxmock.NewResult(1, 1))
			},
			wantErr: true,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			json, err := json.Marshal(testCase.inputOrder)

			if err != nil {
				t.Errorf("error while marshaling order: %s", err)
			}

			testCase.mock(testCase.inputOrder, json)
			err = testCase.s.SaveOrderInRepository(testCase.inputOrder)

			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}

}
