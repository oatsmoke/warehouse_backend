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

func truncateCategories(t *testing.T, testDB *pgxpool.Pool) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const query = `
		TRUNCATE categories
		RESTART IDENTITY CASCADE;`

	if _, err := testDB.Exec(ctx, query); err != nil {
		t.Fatalf("failed to truncate category: %v", err)
	}
}

func addTestCategory(t *testing.T, testDB *pgxpool.Pool) *model.Category {
	t.Helper()
	c := new(model.Category)

	const query = `
		INSERT INTO categories (title) 
		VALUES ($1) 
		RETURNING id, title;`

	if err := testDB.QueryRow(t.Context(), query, generate.RandString(10)).
		Scan(&c.ID, &c.Title); err != nil {
		t.Fatalf("failed to insert test category: %v", err)
	}

	return c
}

func addTestDeletedCategory(t *testing.T, testDB *pgxpool.Pool) *model.Category {
	t.Helper()
	c := new(model.Category)

	const query = `
		INSERT INTO categories (title, deleted_at) 
		VALUES ($1, now()) 
		RETURNING id, title, deleted_at;`

	if err := testDB.QueryRow(t.Context(), query, generate.RandString(10)).
		Scan(&c.ID, &c.Title, &c.DeletedAt); err != nil {
		t.Fatalf("failed to insert test category: %v", err)
	}

	return c
}

func TestNewCategoryRepository(t *testing.T) {
	testDB := postgresql.ConnectTest()
	defer testDB.Close()
	q := queries.New(testDB)

	type args struct {
		queries queries.Querier
	}
	tests := []struct {
		name string
		args args
		want *CategoryRepository
	}{
		{
			name: "create category repository",
			args: args{
				queries: q,
			},
			want: NewCategoryRepository(q),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCategoryRepository(tt.args.queries); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCategoryRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCategoryRepository_Create(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateCategories(t, testDB)
		testDB.Close()
	})
	q := queries.New(testDB)

	type fields struct {
		queries queries.Querier
	}
	type args struct {
		ctx      context.Context
		category *model.Category
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "create category",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				category: &model.Category{
					Title: "test category",
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "create duplicate category",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				category: &model.Category{
					Title: "test category",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CategoryRepository{
				queries: tt.fields.queries,
			}
			got, err := r.Create(tt.args.ctx, tt.args.category)
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

func TestCategoryRepository_Read(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateCategories(t, testDB)
		testDB.Close()
	})
	q := queries.New(testDB)
	c := addTestCategory(t, testDB)

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
		want    *model.Category
		wantErr bool
	}{
		{
			name: "read category",
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
			name: "read non-existing category",
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
			r := &CategoryRepository{
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

func TestCategoryRepository_Update(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateCategories(t, testDB)
		testDB.Close()
	})
	q := queries.New(testDB)
	c := addTestCategory(t, testDB)

	type fields struct {
		queries queries.Querier
	}
	type args struct {
		ctx      context.Context
		category *model.Category
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "update category",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				category: &model.Category{
					ID:    c.ID,
					Title: "update category",
				},
			},
			wantErr: false,
		},
		{
			name: "update non-existing category",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				category: &model.Category{
					ID:    999,
					Title: "update non-existing category",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CategoryRepository{
				queries: tt.fields.queries,
			}
			if err := r.Update(tt.args.ctx, tt.args.category); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCategoryRepository_Delete(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateCategories(t, testDB)
		testDB.Close()
	})
	q := queries.New(testDB)
	c := addTestCategory(t, testDB)

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
			name: "delete category",
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
			name: "delete non-existing category",
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
			r := &CategoryRepository{
				queries: tt.fields.queries,
			}
			if err := r.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCategoryRepository_Restore(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateCategories(t, testDB)
		testDB.Close()
	})
	q := queries.New(testDB)
	c := addTestCategory(t, testDB)
	dc := addTestDeletedCategory(t, testDB)

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
			name: "restore category",
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
			name: "restore non-existing category",
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
			name: "restore not deleted category",
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
			r := &CategoryRepository{
				queries: tt.fields.queries,
			}
			if err := r.Restore(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Restore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCategoryRepository_List(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateCategories(t, testDB)
		testDB.Close()
	})
	q := queries.New(testDB)
	c := addTestCategory(t, testDB)
	dc := addTestDeletedCategory(t, testDB)

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
		want    []*model.Category
		want1   int64
		wantErr bool
	}{
		{
			name: "list categories without deleted",
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
			want:    []*model.Category{c},
			want1:   1,
			wantErr: false,
		},
		{
			name: "list categories with deleted",
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
			want:    []*model.Category{c, dc},
			want1:   2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CategoryRepository{
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
