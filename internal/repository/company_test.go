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

func truncateCompanies(t *testing.T, testDB *pgxpool.Pool) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const query = `
		TRUNCATE companies
		RESTART IDENTITY CASCADE;`

	if _, err := testDB.Exec(ctx, query); err != nil {
		t.Fatalf("failed to truncate company: %v", err)
	}
}

func addTestCompany(t *testing.T, testDB *pgxpool.Pool) *model.Company {
	t.Helper()
	c := new(model.Company)

	const query = `
		INSERT INTO companies (title) 
		VALUES ($1) 
		RETURNING id, title;`

	if err := testDB.QueryRow(t.Context(), query, generate.RandString(10)).
		Scan(&c.ID, &c.Title); err != nil {
		t.Fatalf("failed to insert test company: %v", err)
	}

	return c
}

func addTestDeletedCompany(t *testing.T, testDB *pgxpool.Pool) *model.Company {
	t.Helper()
	c := new(model.Company)

	const query = `
		INSERT INTO companies (title, deleted_at) 
		VALUES ($1, now()) 
		RETURNING id, title, deleted_at;`

	if err := testDB.QueryRow(t.Context(), query, generate.RandString(10)).
		Scan(&c.ID, &c.Title, &c.DeletedAt); err != nil {
		t.Fatalf("failed to insert test company: %v", err)
	}

	return c
}

func TestNewCompanyRepository(t *testing.T) {
	testDB := postgresql.ConnectTest()
	defer testDB.Close()
	q := queries.New(testDB)

	type args struct {
		queries queries.Querier
	}
	tests := []struct {
		name string
		args args
		want *CompanyRepository
	}{
		{
			name: "create company repository",
			args: args{
				queries: q,
			},
			want: NewCompanyRepository(q),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCompanyRepository(tt.args.queries); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCompanyRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompanyRepository_Create(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateCompanies(t, testDB)
		testDB.Close()
	})
	q := queries.New(testDB)

	type fields struct {
		queries queries.Querier
	}
	type args struct {
		ctx     context.Context
		company *model.Company
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "create company",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				company: &model.Company{
					Title: "test company",
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "create duplicate company",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				company: &model.Company{
					Title: "test company",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CompanyRepository{
				queries: tt.fields.queries,
			}
			got, err := r.Create(tt.args.ctx, tt.args.company)
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

func TestCompanyRepository_Read(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateCompanies(t, testDB)
		testDB.Close()
	})
	q := queries.New(testDB)
	c := addTestCompany(t, testDB)

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
		want    *model.Company
		wantErr bool
	}{
		{
			name: "read company",
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
			name: "read non-existing company",
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
			r := &CompanyRepository{
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

func TestCompanyRepository_Update(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateCompanies(t, testDB)
		testDB.Close()
	})
	q := queries.New(testDB)
	c := addTestCompany(t, testDB)

	type fields struct {
		queries queries.Querier
	}
	type args struct {
		ctx     context.Context
		company *model.Company
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "update company",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				company: &model.Company{
					ID:    c.ID,
					Title: "update company",
				},
			},
			wantErr: false,
		},
		{
			name: "update non-existing company",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				company: &model.Company{
					ID:    999,
					Title: "update non-existing company",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CompanyRepository{
				queries: tt.fields.queries,
			}
			if err := r.Update(tt.args.ctx, tt.args.company); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCompanyRepository_Delete(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateCompanies(t, testDB)
		testDB.Close()
	})
	q := queries.New(testDB)
	c := addTestCompany(t, testDB)

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
			name: "delete company",
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
			name: "delete non-existing company",
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
			r := &CompanyRepository{
				queries: tt.fields.queries,
			}
			if err := r.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCompanyRepository_Restore(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateCompanies(t, testDB)
		testDB.Close()
	})
	q := queries.New(testDB)
	c := addTestCompany(t, testDB)
	dc := addTestDeletedCompany(t, testDB)

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
			name: "restore company",
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
			name: "restore non-existing company",
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
			name: "restore not deleted company",
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
			r := &CompanyRepository{
				queries: tt.fields.queries,
			}
			if err := r.Restore(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Restore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCompanyRepository_List(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateCompanies(t, testDB)
		testDB.Close()
	})
	q := queries.New(testDB)
	c := addTestCompany(t, testDB)
	dc := addTestDeletedCompany(t, testDB)

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
		want    []*model.Company
		want1   int64
		wantErr bool
	}{
		{
			name: "list companies without deleted",
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
			want:    []*model.Company{c},
			want1:   1,
			wantErr: false,
		},
		{
			name: "list companies with deleted",
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
			want:    []*model.Company{c, dc},
			want1:   2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CompanyRepository{
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
