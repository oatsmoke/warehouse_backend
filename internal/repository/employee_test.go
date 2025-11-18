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

func truncateEmployees(t *testing.T, testDB *pgxpool.Pool) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const query = `
		TRUNCATE employees, departments
		RESTART IDENTITY CASCADE;`

	if _, err := testDB.Exec(ctx, query); err != nil {
		t.Fatalf("failed to truncate employee: %v", err)
	}
}

func addTestEmployee(t *testing.T, testDB *pgxpool.Pool) *model.Employee {
	t.Helper()
	e := new(model.Employee)

	const query = `
		INSERT INTO employees (last_name, first_name, middle_name, phone)
		VALUES ($1, $2, $3, $4)
		RETURNING id, last_name, first_name, middle_name, phone;`

	if err := testDB.QueryRow(t.Context(), query, generate.RandString(10), generate.RandString(10), generate.RandString(10), generate.RandString(10)).
		Scan(&e.ID, &e.LastName, &e.FirstName, &e.MiddleName, &e.Phone); err != nil {
		t.Fatalf("failed to insert test employee: %v", err)
	}

	e.Department = &model.Department{}

	return e
}

func addTestDeletedEmployee(t *testing.T, testDB *pgxpool.Pool) *model.Employee {
	t.Helper()
	e := new(model.Employee)

	const query = `
		INSERT INTO employees (last_name, first_name, middle_name, phone, deleted_at)
		VALUES ($1, $2, $3, $4, now())
		RETURNING id, last_name, first_name, middle_name, phone, deleted_at;`

	if err := testDB.QueryRow(t.Context(), query, generate.RandString(10), generate.RandString(10), generate.RandString(10), generate.RandString(10)).
		Scan(&e.ID, &e.LastName, &e.FirstName, &e.MiddleName, &e.Phone, &e.DeletedAt); err != nil {
		t.Fatalf("failed to insert test profile: %v", err)
	}

	e.Department = &model.Department{}

	return e
}

func TestNewEmployeeRepository(t *testing.T) {
	testDB := postgresql.ConnectTest()
	defer testDB.Close()
	q := queries.New(testDB)

	type args struct {
		queries queries.Querier
	}
	tests := []struct {
		name string
		args args
		want *EmployeeRepository
	}{
		{
			name: "create employee repository",
			args: args{
				queries: q,
			},
			want: NewEmployeeRepository(q),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEmployeeRepository(tt.args.queries); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEmployeeRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmployeeRepository_Create(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateEmployees(t, testDB)
		testDB.Close()
	})
	q := queries.New(testDB)

	type fields struct {
		queries queries.Querier
	}
	type args struct {
		ctx      context.Context
		employee *model.Employee
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "create employee",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				employee: &model.Employee{
					LastName:   "test",
					FirstName:  "test",
					MiddleName: "test",
					Phone:      "1234567890",
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "create duplicate employee",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				employee: &model.Employee{
					LastName:   "test",
					FirstName:  "test",
					MiddleName: "test",
					Phone:      "1234567890",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &EmployeeRepository{
				queries: tt.fields.queries,
			}
			got, err := r.Create(tt.args.ctx, tt.args.employee)
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

func TestEmployeeRepository_Read(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateEmployees(t, testDB)
		testDB.Close()
	})
	q := queries.New(testDB)
	e := addTestEmployee(t, testDB)

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
		want    *model.Employee
		wantErr bool
	}{
		{
			name: "read employee",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				id:  e.ID,
			},
			want:    e,
			wantErr: false,
		},
		{
			name: "read non-existing employee",
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
			r := &EmployeeRepository{
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

func TestEmployeeRepository_Update(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateEmployees(t, testDB)
		testDB.Close()
	})
	q := queries.New(testDB)
	e := addTestEmployee(t, testDB)

	type fields struct {
		queries queries.Querier
	}
	type args struct {
		ctx      context.Context
		employee *model.Employee
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "update employee",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				employee: &model.Employee{
					ID:         e.ID,
					LastName:   "update test",
					FirstName:  "test",
					MiddleName: "test",
					Phone:      "1234567890",
				},
			},
			wantErr: false,
		},
		{
			name: "update non-existing employee",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				employee: &model.Employee{
					ID:         999,
					LastName:   "update non-existing test",
					FirstName:  "test",
					MiddleName: "test",
					Phone:      "1234567890",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &EmployeeRepository{
				queries: tt.fields.queries,
			}
			if err := r.Update(tt.args.ctx, tt.args.employee); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEmployeeRepository_Delete(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateEmployees(t, testDB)
		testDB.Close()
	})
	q := queries.New(testDB)
	e := addTestEmployee(t, testDB)

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
			name: "delete employee",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				id:  e.ID,
			},
			wantErr: false,
		},
		{
			name: "delete non-existing employee",
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
			r := &EmployeeRepository{
				queries: tt.fields.queries,
			}
			if err := r.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEmployeeRepository_Restore(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateEmployees(t, testDB)
		testDB.Close()
	})
	q := queries.New(testDB)
	e := addTestEmployee(t, testDB)
	de := addTestDeletedEmployee(t, testDB)

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
			name: "restore employee",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				id:  de.ID,
			},
			wantErr: false,
		},
		{
			name: "restore non-existing employee",
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
			name: "restore not deleted employee",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				id:  e.ID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &EmployeeRepository{
				queries: tt.fields.queries,
			}
			if err := r.Restore(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Restore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEmployeeRepository_List(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateEmployees(t, testDB)
		testDB.Close()
	})
	q := queries.New(testDB)
	e := addTestEmployee(t, testDB)
	de := addTestDeletedEmployee(t, testDB)

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
		want    []*model.Employee
		want1   int64
		wantErr bool
	}{
		{
			name: "list employees without deleted",
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
			want:    []*model.Employee{e},
			want1:   1,
			wantErr: false,
		},
		{
			name: "list employees with deleted",
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
			want:    []*model.Employee{e, de},
			want1:   2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &EmployeeRepository{
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

func TestEmployeeRepository_SetDepartment(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateEmployees(t, testDB)
		testDB.Close()
	})
	q := queries.New(testDB)
	d := addTestDepartment(t, testDB)
	e := addTestEmployee(t, testDB)

	type fields struct {
		queries queries.Querier
	}
	type args struct {
		ctx          context.Context
		id           int64
		departmentID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "set department",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx:          t.Context(),
				id:           e.ID,
				departmentID: d.ID,
			},
			wantErr: false,
		},
		{
			name: "set non-existing department",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx:          t.Context(),
				id:           e.ID,
				departmentID: 999,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &EmployeeRepository{
				queries: tt.fields.queries,
			}
			if err := r.SetDepartment(tt.args.ctx, tt.args.id, tt.args.departmentID); (err != nil) != tt.wantErr {
				t.Errorf("SetDepartment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
