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

func truncateDepartments(t *testing.T, testDB *pgxpool.Pool) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const query = `
		TRUNCATE departments
		RESTART IDENTITY CASCADE;`

	if _, err := testDB.Exec(ctx, query); err != nil {
		t.Fatalf("failed to truncate department: %v", err)
	}
}

func addTestDepartment(t *testing.T, testDB *pgxpool.Pool) *model.Department {
	t.Helper()
	d := new(model.Department)

	const query = `
		INSERT INTO departments (title) 
		VALUES ($1) 
		RETURNING id, title;`

	if err := testDB.QueryRow(t.Context(), query, generate.RandString(10)).
		Scan(&d.ID, &d.Title); err != nil {
		t.Fatalf("failed to insert test department: %v", err)
	}

	return d
}

func addTestDeletedDepartment(t *testing.T, testDB *pgxpool.Pool) *model.Department {
	t.Helper()
	d := new(model.Department)

	const query = `
		INSERT INTO departments (title, deleted_at) 
		VALUES ($1, now()) 
		RETURNING id, title, deleted_at;`

	if err := testDB.QueryRow(t.Context(), query, generate.RandString(10)).
		Scan(&d.ID, &d.Title, &d.DeletedAt); err != nil {
		t.Fatalf("failed to insert test department: %v", err)
	}

	return d
}
func TestNewDepartmentRepository(t *testing.T) {
	testDB := postgresql.ConnectTest()
	defer testDB.Close()
	q := queries.New(testDB)

	type args struct {
		queries queries.Querier
	}
	tests := []struct {
		name string
		args args
		want *DepartmentRepository
	}{
		{
			name: "create department repository",
			args: args{
				queries: q,
			},
			want: NewDepartmentRepository(q),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDepartmentRepository(tt.args.queries); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDepartmentRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDepartmentRepository_Create(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateDepartments(t, testDB)
		testDB.Close()
	})
	q := queries.New(testDB)

	type fields struct {
		queries queries.Querier
	}
	type args struct {
		ctx        context.Context
		department *model.Department
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "create department",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				department: &model.Department{
					Title: "test department",
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "create duplicate department",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				department: &model.Department{
					Title: "test department",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &DepartmentRepository{
				queries: tt.fields.queries,
			}
			got, err := r.Create(tt.args.ctx, tt.args.department)
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

func TestDepartmentRepository_Read(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateDepartments(t, testDB)
		testDB.Close()
	})
	q := queries.New(testDB)
	d := addTestDepartment(t, testDB)

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
		want    *model.Department
		wantErr bool
	}{
		{
			name: "read department",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				id:  d.ID,
			},
			want:    d,
			wantErr: false,
		},
		{
			name: "read non-existing department",
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
			r := &DepartmentRepository{
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

func TestDepartmentRepository_Update(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateDepartments(t, testDB)
		testDB.Close()
	})
	q := queries.New(testDB)
	d := addTestDepartment(t, testDB)

	type fields struct {
		queries queries.Querier
	}
	type args struct {
		ctx        context.Context
		department *model.Department
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "update department",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				department: &model.Department{
					ID:    d.ID,
					Title: "update department",
				},
			},
			wantErr: false,
		},
		{
			name: "update non-existing department",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				department: &model.Department{
					ID:    999,
					Title: "update non-existing department",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &DepartmentRepository{
				queries: tt.fields.queries,
			}
			if err := r.Update(tt.args.ctx, tt.args.department); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDepartmentRepository_Delete(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateDepartments(t, testDB)
		testDB.Close()
	})
	q := queries.New(testDB)
	d := addTestDepartment(t, testDB)

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
			name: "delete department",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				id:  d.ID,
			},
			wantErr: false,
		},
		{
			name: "delete non-existing department",
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
			r := &DepartmentRepository{
				queries: tt.fields.queries,
			}
			if err := r.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDepartmentRepository_Restore(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateDepartments(t, testDB)
		testDB.Close()
	})
	q := queries.New(testDB)
	d := addTestDepartment(t, testDB)
	dd := addTestDeletedDepartment(t, testDB)

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
			name: "restore department",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				id:  dd.ID,
			},
			wantErr: false,
		},
		{
			name: "restore non-existing department",
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
			name: "restore not deleted department",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				id:  d.ID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &DepartmentRepository{
				queries: tt.fields.queries,
			}
			if err := r.Restore(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Restore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDepartmentRepository_List(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateDepartments(t, testDB)
		testDB.Close()
	})
	q := queries.New(testDB)
	d := addTestDepartment(t, testDB)
	dd := addTestDeletedDepartment(t, testDB)

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
		want    []*model.Department
		want1   int64
		wantErr bool
	}{
		{
			name: "list departments without deleted",
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
			want:    []*model.Department{d},
			want1:   1,
			wantErr: false,
		},
		{
			name: "list departments with deleted",
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
			want:    []*model.Department{d, dd},
			want1:   2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &DepartmentRepository{
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
