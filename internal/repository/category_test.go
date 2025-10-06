package repository

import (
	"context"
	"log"
	"reflect"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oatsmoke/warehouse_backend/internal/dto"
	"github.com/oatsmoke/warehouse_backend/internal/lib/postgresql"
	"github.com/oatsmoke/warehouse_backend/internal/model"
)

func truncateCategories() {
	_, err := postgresql.ConnectTest().Exec(context.Background(), "TRUNCATE categories RESTART IDENTITY CASCADE;")
	if err != nil {
		log.Fatal(err)
	}
}

func addTestCategory(ctx context.Context) *model.Category {
	c := new(model.Category)

	const query = `
		INSERT INTO categories (title) 
		VALUES ($1) 
		RETURNING id,title;`

	if err := postgresql.ConnectTest().QueryRow(ctx, query, "test_category").
		Scan(&c.ID, &c.Title); err != nil {
		log.Fatalf("failed to insert test category: %v\n", err)
	}

	return c
}

func addTestDeletedCategory(ctx context.Context) *model.Category {
	c := new(model.Category)

	const query = `
		INSERT INTO categories (title, deleted_at) 
		VALUES ($1, now()) 
		RETURNING id,title,deleted_at;`

	if err := postgresql.ConnectTest().QueryRow(ctx, query, "delete_category").
		Scan(&c.ID, &c.Title, &c.DeletedAt); err != nil {
		log.Fatalf("failed to insert test category: %v\n", err)
	}

	return c
}

func TestNewCategoryRepository(t *testing.T) {
	type args struct {
		postgresDB *pgxpool.Pool
	}
	tests := []struct {
		name string
		args args
		want *CategoryRepository
	}{
		{
			name: "create category repository",
			args: args{
				postgresDB: postgresql.ConnectTest(),
			},
			want: &CategoryRepository{
				postgresDB: postgresql.ConnectTest(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCategoryRepository(tt.args.postgresDB); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCategoryRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCategoryRepository_Create(t *testing.T) {
	defer truncateCategories()

	type fields struct {
		postgresDB *pgxpool.Pool
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
				postgresDB: postgresql.ConnectTest(),
			},
			args: args{
				ctx: context.Background(),
				category: &model.Category{
					Title: "test_category",
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "create duplicate category",
			fields: fields{
				postgresDB: postgresql.ConnectTest(),
			},
			args: args{
				ctx: context.Background(),
				category: &model.Category{
					Title: "test_category",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CategoryRepository{
				postgresDB: tt.fields.postgresDB,
			}
			got, err := r.Create(tt.args.ctx, tt.args.category)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCategoryRepository_Read(t *testing.T) {
	defer truncateCategories()

	c := addTestCategory(t.Context())

	type fields struct {
		postgresDB *pgxpool.Pool
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
				postgresDB: postgresql.ConnectTest(),
			},
			args: args{
				ctx: context.Background(),
				id:  c.ID,
			},
			want: &model.Category{
				ID:    c.ID,
				Title: c.Title,
			},
			wantErr: false,
		},
		{
			name: "read non-existing category",
			fields: fields{
				postgresDB: postgresql.ConnectTest(),
			},
			args: args{
				ctx: context.Background(),
				id:  999,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CategoryRepository{
				postgresDB: tt.fields.postgresDB,
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
	defer truncateCategories()

	c := addTestCategory(t.Context())

	type fields struct {
		postgresDB *pgxpool.Pool
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
				postgresDB: postgresql.ConnectTest(),
			},
			args: args{
				ctx: context.Background(),
				category: &model.Category{
					ID:    c.ID,
					Title: "updated_category",
				},
			},
			wantErr: false,
		},
		{
			name: "update non-existing category",
			fields: fields{
				postgresDB: postgresql.ConnectTest(),
			},
			args: args{
				ctx: context.Background(),
				category: &model.Category{
					ID:    999,
					Title: "non_existing_category",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CategoryRepository{
				postgresDB: tt.fields.postgresDB,
			}
			if err := r.Update(tt.args.ctx, tt.args.category); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCategoryRepository_Delete(t *testing.T) {
	defer truncateCategories()

	c := addTestCategory(t.Context())

	type fields struct {
		postgresDB *pgxpool.Pool
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
				postgresDB: postgresql.ConnectTest(),
			},
			args: args{
				ctx: context.Background(),
				id:  c.ID,
			},
			wantErr: false,
		},
		{
			name: "delete non-existing category",
			fields: fields{
				postgresDB: postgresql.ConnectTest(),
			},
			args: args{
				ctx: context.Background(),
				id:  999,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CategoryRepository{
				postgresDB: tt.fields.postgresDB,
			}
			if err := r.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCategoryRepository_List(t *testing.T) {
	defer truncateCategories()

	c := addTestCategory(t.Context())
	dc := addTestDeletedCategory(t.Context())

	type fields struct {
		postgresDB *pgxpool.Pool
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
		wantErr bool
	}{
		{
			name: "list categories without deleted",
			fields: fields{
				postgresDB: postgresql.ConnectTest(),
			},
			args: args{
				ctx: context.Background(),
				qp: &dto.QueryParams{
					WithDeleted: false,
				},
			},
			want: []*model.Category{
				{
					ID:    c.ID,
					Title: c.Title,
				},
			},
			wantErr: false,
		},
		{
			name: "list categories with deleted",
			fields: fields{
				postgresDB: postgresql.ConnectTest(),
			},
			args: args{
				ctx: context.Background(),
				qp: &dto.QueryParams{
					WithDeleted: false,
				},
			},
			want: []*model.Category{
				{
					ID:        dc.ID,
					Title:     dc.Title,
					DeletedAt: dc.DeletedAt,
				},
				{
					ID:    c.ID,
					Title: c.Title,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CategoryRepository{
				postgresDB: tt.fields.postgresDB,
			}
			got, err := r.List(tt.args.ctx, tt.args.qp)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCategoryRepository_Restore(t *testing.T) {
	defer truncateCategories()

	c := addTestCategory(t.Context())
	dc := addTestDeletedCategory(t.Context())

	type fields struct {
		postgresDB *pgxpool.Pool
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
				postgresDB: postgresql.ConnectTest(),
			},
			args: args{
				ctx: context.Background(),
				id:  dc.ID,
			},
			wantErr: false,
		},
		{
			name: "restore non-existing category",
			fields: fields{
				postgresDB: postgresql.ConnectTest(),
			},
			args: args{
				ctx: context.Background(),
				id:  999,
			},
			wantErr: true,
		},
		{
			name: "restore not deleted category",
			fields: fields{
				postgresDB: postgresql.ConnectTest(),
			},
			args: args{
				ctx: context.Background(),
				id:  c.ID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CategoryRepository{
				postgresDB: tt.fields.postgresDB,
			}
			if err := r.Restore(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Restore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
