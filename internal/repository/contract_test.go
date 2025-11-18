package repository

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	queries "github.com/oatsmoke/warehouse_backend/internal/db"
	"github.com/oatsmoke/warehouse_backend/internal/dto"
	"github.com/oatsmoke/warehouse_backend/internal/lib/generate"
	"github.com/oatsmoke/warehouse_backend/internal/lib/postgresql"
	"github.com/oatsmoke/warehouse_backend/internal/model"
)

func truncateContracts(t *testing.T, testDB *pgxpool.Pool) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const query = `
		TRUNCATE contracts
		RESTART IDENTITY CASCADE;`

	if _, err := testDB.Exec(ctx, query); err != nil {
		t.Fatalf("failed to truncate contract: %v", err)
	}
}

func addTestContract(t *testing.T, testDB *pgxpool.Pool) *model.Contract {
	t.Helper()
	c := new(model.Contract)

	const query = `
		INSERT INTO contracts (number, address) 
		VALUES ($1, $2) 
		RETURNING id, number, address;`

	if err := testDB.QueryRow(t.Context(), query, generate.RandString(10), generate.RandString(10)).
		Scan(&c.ID, &c.Number, &c.Address); err != nil {
		t.Fatalf("failed to insert test contract: %v", err)
	}

	return c
}

func addTestDeletedContract(t *testing.T, testDB *pgxpool.Pool) *model.Contract {
	t.Helper()
	c := new(model.Contract)

	const query = `
		INSERT INTO contracts (number, address, deleted_at) 
		VALUES ($1, $2, now()) 
		RETURNING id, number, address, deleted_at;`

	if err := testDB.QueryRow(t.Context(), query, generate.RandString(10), generate.RandString(10)).
		Scan(&c.ID, &c.Number, &c.Address, &c.DeletedAt); err != nil {
		t.Fatalf("failed to insert test contract: %v", err)
	}

	return c
}

func TestNewContractRepository(t *testing.T) {
	testDB := postgresql.ConnectTest()
	defer testDB.Close()
	q := queries.New(testDB)

	type args struct {
		queries queries.Querier
	}
	tests := []struct {
		name string
		args args
		want *ContractRepository
	}{
		{
			name: "create contract repository",
			args: args{
				queries: q,
			},
			want: NewContractRepository(q),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewContractRepository(tt.args.queries); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewContractRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContractRepository_Create(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateContracts(t, testDB)
		testDB.Close()
	})
	q := queries.New(testDB)

	type fields struct {
		queries queries.Querier
	}
	type args struct {
		ctx      context.Context
		contract *model.Contract
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "create contract",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				contract: &model.Contract{
					Number:  "800000",
					Address: "test contract",
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "create duplicate contract",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				contract: &model.Contract{
					Number:  "800000",
					Address: "test contract",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ContractRepository{
				queries: tt.fields.queries,
			}
			got, err := r.Create(tt.args.ctx, tt.args.contract)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContractRepository_Read(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateContracts(t, testDB)
		testDB.Close()
	})
	q := queries.New(testDB)
	c := addTestContract(t, testDB)

	type fields struct {
		queries queries.Querier
	}
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Contract
		wantErr bool
	}{
		{
			name: "read contract",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				id:  c.ID,
			},
			want:    c,
			wantErr: false,
		},
		{
			name: "read non-existing contract",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				id:  999,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ContractRepository{
				queries: tt.fields.queries,
			}
			got, err := r.Read(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Read() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContractRepository_Update(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateContracts(t, testDB)
		testDB.Close()
	})
	q := queries.New(testDB)
	c := addTestContract(t, testDB)

	type fields struct {
		queries queries.Querier
	}
	type args struct {
		ctx      context.Context
		contract *model.Contract
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "update contract",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				contract: &model.Contract{
					ID:      c.ID,
					Number:  "800000",
					Address: "update contract",
				},
			},
			wantErr: false,
		},
		{
			name: "update non-existing contract",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				contract: &model.Contract{
					ID:      999,
					Number:  "800000",
					Address: "update non-existing contract",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ContractRepository{
				queries: tt.fields.queries,
			}
			if err := r.Update(tt.args.ctx, tt.args.contract); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContractRepository_Delete(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateContracts(t, testDB)
		testDB.Close()
	})
	q := queries.New(testDB)
	c := addTestContract(t, testDB)

	type fields struct {
		queries queries.Querier
	}
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "delete contract",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				id:  c.ID,
			},
			wantErr: false,
		},
		{
			name: "delete non-existing contract",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				id:  999,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ContractRepository{
				queries: tt.fields.queries,
			}
			if err := r.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContractRepository_Restore(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateContracts(t, testDB)
		testDB.Close()
	})
	q := queries.New(testDB)
	c := addTestContract(t, testDB)
	dc := addTestDeletedContract(t, testDB)

	type fields struct {
		queries queries.Querier
	}
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "restore contract",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				id:  dc.ID,
			},
			wantErr: false,
		},
		{
			name: "restore non-existing contract",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				id:  999,
			},
			wantErr: true,
		},
		{
			name: "restore not deleted contract",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				id:  c.ID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ContractRepository{
				queries: tt.fields.queries,
			}
			if err := r.Restore(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Restore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContractRepository_List(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateContracts(t, testDB)
		testDB.Close()
	})
	q := queries.New(testDB)
	c := addTestContract(t, testDB)
	dc := addTestDeletedContract(t, testDB)

	type fields struct {
		queries queries.Querier
	}
	type args struct {
		ctx context.Context
		qp  *dto.QueryParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*model.Contract
		want1   int64
		wantErr bool
	}{
		{
			name: "list contracts without deleted",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				qp: &dto.QueryParams{
					SortColumn:       "id",
					SortOrder:        "asc",
					PaginationLimit:  50,
					PaginationOffset: 0,
				},
			},
			want:    []*model.Contract{c},
			want1:   1,
			wantErr: false,
		},
		{
			name: "list contracts with deleted",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				qp: &dto.QueryParams{
					WithDeleted:      true,
					SortColumn:       "id",
					SortOrder:        "asc",
					PaginationLimit:  50,
					PaginationOffset: 0,
				},
			},
			want:    []*model.Contract{c, dc},
			want1:   2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ContractRepository{
				queries: tt.fields.queries,
			}
			got, got1, err := r.List(tt.args.ctx, tt.args.qp)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("List() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
