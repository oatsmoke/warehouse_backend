package repository

import (
	"context"
	"reflect"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oatsmoke/warehouse_backend/internal/lib/env"
	"github.com/oatsmoke/warehouse_backend/internal/lib/postgresql"
	"github.com/oatsmoke/warehouse_backend/internal/model"
)

var conn = postgresql.Connect(context.Background(), env.GetTestPostgresDsn())

func truncateCategories() {
	_, err := conn.Exec(context.Background(), "TRUNCATE categories RESTART IDENTITY CASCADE;")
	if err != nil {
		panic(err)
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
				postgresDB: conn,
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
				postgresDB: conn,
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

func TestCategoryRepository_Delete(t *testing.T) {
	defer truncateCategories()

	var id int64
	if err := conn.QueryRow(context.Background(), "INSERT INTO categories (title) VALUES ('test_category') RETURNING id;").
		Scan(&id); err != nil {
		t.Fatalf("failed to insert test category: %v", err)
	}

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
				postgresDB: conn,
			},
			args: args{
				ctx: context.Background(),
				id:  id,
			},
			wantErr: false,
		},
		{
			name: "delete non-existing category",
			fields: fields{
				postgresDB: conn,
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

	ndc := model.Category{}
	if err := conn.QueryRow(context.Background(), "INSERT INTO categories (title) VALUES ('not_deleted_category') RETURNING id,title;").
		Scan(&ndc.ID, &ndc.Title); err != nil {
		t.Fatalf("failed to insert test category: %v", err)
	}

	dc := model.Category{}
	if err := conn.QueryRow(context.Background(), "INSERT INTO categories (title, deleted_at) VALUES ('delete_category',now()) RETURNING id,title,deleted_at;").
		Scan(&dc.ID, &dc.Title, &dc.DeletedAt); err != nil {
		t.Fatalf("failed to insert test category: %v", err)
	}

	type fields struct {
		postgresDB *pgxpool.Pool
	}
	type args struct {
		ctx         context.Context
		withDeleted bool
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
				postgresDB: conn,
			},
			args: args{
				ctx:         context.Background(),
				withDeleted: false,
			},
			want: []*model.Category{
				{
					ID:    ndc.ID,
					Title: ndc.Title,
				},
			},
			wantErr: false,
		},
		{
			name: "list categories with deleted",
			fields: fields{
				postgresDB: conn,
			},
			args: args{
				ctx:         context.Background(),
				withDeleted: true,
			},
			want: []*model.Category{
				{
					ID:        dc.ID,
					Title:     dc.Title,
					DeletedAt: dc.DeletedAt,
				},
				{
					ID:    ndc.ID,
					Title: ndc.Title,
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
			got, err := r.List(tt.args.ctx, tt.args.withDeleted)
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

func TestCategoryRepository_Read(t *testing.T) {
	defer truncateCategories()

	c := model.Category{}
	if err := conn.QueryRow(context.Background(), "INSERT INTO categories (title) VALUES ('test_category') RETURNING id,title;").
		Scan(&c.ID, &c.Title); err != nil {
		t.Fatalf("failed to insert test category: %v", err)
	}

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
			name: "read existing category",
			fields: fields{
				postgresDB: conn,
			},
			args: args{
				ctx: context.Background(),
				id:  1,
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
				postgresDB: conn,
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

func TestCategoryRepository_Restore(t *testing.T) {
	defer truncateCategories()

	c := model.Category{}
	if err := conn.QueryRow(context.Background(), "INSERT INTO categories (title, deleted_at) VALUES ('test_category',now()) RETURNING id,title,deleted_at;").
		Scan(&c.ID, &c.Title, &c.DeletedAt); err != nil {
		t.Fatalf("failed to insert test category: %v", err)
	}

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
				postgresDB: conn,
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			wantErr: false,
		},
		{
			name: "restore non-existing category",
			fields: fields{
				postgresDB: conn,
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
				postgresDB: conn,
			},
			args: args{
				ctx: context.Background(),
				id:  1,
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

func TestCategoryRepository_Update(t *testing.T) {
	defer truncateCategories()

	var id int64
	if err := conn.QueryRow(context.Background(), "INSERT INTO categories (title) VALUES ('test_category') RETURNING id;").
		Scan(&id); err != nil {
		t.Fatalf("failed to insert test category: %v", err)
	}

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
				postgresDB: conn,
			},
			args: args{
				ctx: context.Background(),
				category: &model.Category{
					ID:    id,
					Title: "updated_category",
				},
			},
			wantErr: false,
		},
		{
			name: "update non-existing category",
			fields: fields{
				postgresDB: conn,
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
			name: "create new category repository",
			args: args{
				postgresDB: conn,
			},
			want: &CategoryRepository{
				postgresDB: conn,
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
