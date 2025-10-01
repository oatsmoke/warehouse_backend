package repository

import (
	"context"
	"log"
	"reflect"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oatsmoke/warehouse_backend/internal/model"
)

func truncateProfiles() {
	_, err := testConn.Exec(context.Background(), "TRUNCATE profiles RESTART IDENTITY CASCADE;")
	if err != nil {
		log.Fatal(err)
	}
}

func addTestProfile(ctx context.Context, categoryID int64) *model.Profile {
	p := new(model.Profile)

	const query = `
		INSERT INTO profiles (title, category)
		VALUES ($1, $2)
		RETURNING id, title;`

	if err := testConn.QueryRow(ctx, query, "test_profile", categoryID).
		Scan(&p.ID, &p.Title); err != nil {
		log.Fatalf("failed to insert test profile: %v\n", err)
	}

	return p
}

func addTestDeletedProfile(ctx context.Context, categoryID int64) *model.Profile {
	p := new(model.Profile)

	const query = `
		INSERT INTO profiles (title, category, deleted_at)
		VALUES ($1, $2, now())
		RETURNING id, title, deleted_at;`

	if err := testConn.QueryRow(ctx, query, "delete_profile", categoryID).
		Scan(&p.ID, &p.Title, &p.DeletedAt); err != nil {
		log.Fatalf("failed to insert test profile: %v\n", err)
	}

	return p
}

func TestNewProfileRepository(t *testing.T) {
	type args struct {
		postgresDB *pgxpool.Pool
	}
	tests := []struct {
		name string
		args args
		want *ProfileRepository
	}{
		{
			name: "create profile repository",
			args: args{
				postgresDB: testConn,
			},
			want: &ProfileRepository{
				postgresDB: testConn,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewProfileRepository(tt.args.postgresDB); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProfileRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProfileRepository_Create(t *testing.T) {
	defer func() {
		truncateProfiles()
		truncateCategories()
	}()

	c := addTestCategory(t.Context())

	type fields struct {
		postgresDB *pgxpool.Pool
	}
	type args struct {
		ctx     context.Context
		profile *model.Profile
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "create profile",
			fields: fields{
				postgresDB: testConn,
			},
			args: args{
				ctx: context.Background(),
				profile: &model.Profile{
					Title: "test_profile",
					Category: &model.Category{
						ID: c.ID,
					},
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "create duplicate profile",
			fields: fields{
				postgresDB: testConn,
			},
			args: args{
				ctx: context.Background(),
				profile: &model.Profile{
					Title: "test_profile",
					Category: &model.Category{
						ID: c.ID,
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ProfileRepository{
				postgresDB: tt.fields.postgresDB,
			}
			got, err := r.Create(tt.args.ctx, tt.args.profile)
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

func TestProfileRepository_Read(t *testing.T) {
	defer func() {
		truncateProfiles()
		truncateCategories()
	}()

	c := addTestCategory(t.Context())
	p := addTestProfile(t.Context(), c.ID)

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
		want    *model.Profile
		wantErr bool
	}{
		{
			name: "read profile",
			fields: fields{
				postgresDB: testConn,
			},
			args: args{
				ctx: context.Background(),
				id:  p.ID,
			},
			want: &model.Profile{
				ID:    p.ID,
				Title: p.Title,
				Category: &model.Category{
					ID:    c.ID,
					Title: c.Title,
				},
			},
			wantErr: false,
		},
		{
			name: "read non-existing profile",
			fields: fields{
				postgresDB: testConn,
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
			r := &ProfileRepository{
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

func TestProfileRepository_Update(t *testing.T) {
	defer func() {
		truncateProfiles()
		truncateCategories()
	}()

	c := addTestCategory(t.Context())
	p := addTestProfile(t.Context(), c.ID)

	type fields struct {
		postgresDB *pgxpool.Pool
	}
	type args struct {
		ctx     context.Context
		profile *model.Profile
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "update profile",
			fields: fields{
				postgresDB: testConn,
			},
			args: args{
				ctx: context.Background(),
				profile: &model.Profile{
					ID:    p.ID,
					Title: "updated_profile",
					Category: &model.Category{
						ID: c.ID,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "update non-existing profile",
			fields: fields{
				postgresDB: testConn,
			},
			args: args{
				ctx: context.Background(),
				profile: &model.Profile{
					ID:    999,
					Title: "non_existing_profile",
					Category: &model.Category{
						ID: c.ID,
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ProfileRepository{
				postgresDB: tt.fields.postgresDB,
			}
			if err := r.Update(tt.args.ctx, tt.args.profile); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProfileRepository_Delete(t *testing.T) {
	defer func() {
		truncateProfiles()
		truncateCategories()
	}()

	c := addTestCategory(t.Context())
	p := addTestProfile(t.Context(), c.ID)

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
			name: "delete profile",
			fields: fields{
				postgresDB: testConn,
			},
			args: args{
				ctx: context.Background(),
				id:  p.ID,
			},
			wantErr: false,
		},
		{
			name: "delete non-existing profile",
			fields: fields{
				postgresDB: testConn,
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
			r := &ProfileRepository{
				postgresDB: tt.fields.postgresDB,
			}
			if err := r.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProfileRepository_List(t *testing.T) {
	defer func() {
		truncateProfiles()
		truncateCategories()
	}()

	c := addTestCategory(t.Context())
	p := addTestProfile(t.Context(), c.ID)
	dp := addTestDeletedProfile(t.Context(), c.ID)

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
		want    []*model.Profile
		wantErr bool
	}{
		{
			name: "list profiles without deleted",
			fields: fields{
				postgresDB: testConn,
			},
			args: args{
				ctx:         context.Background(),
				withDeleted: false,
			},
			want: []*model.Profile{
				{
					ID:    p.ID,
					Title: p.Title,
					Category: &model.Category{
						ID:    c.ID,
						Title: c.Title,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "list profiles with deleted",
			fields: fields{
				postgresDB: testConn,
			},
			args: args{
				ctx:         context.Background(),
				withDeleted: true,
			},
			want: []*model.Profile{
				{
					ID:        dp.ID,
					Title:     dp.Title,
					DeletedAt: dp.DeletedAt,
					Category: &model.Category{
						ID:    c.ID,
						Title: c.Title,
					},
				},
				{
					ID:    p.ID,
					Title: p.Title,
					Category: &model.Category{
						ID:    c.ID,
						Title: c.Title,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ProfileRepository{
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

func TestProfileRepository_Restore(t *testing.T) {
	defer func() {
		truncateProfiles()
		truncateCategories()
	}()

	c := addTestCategory(t.Context())
	p := addTestDeletedProfile(t.Context(), c.ID)

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
			name: "restore profile",
			fields: fields{
				postgresDB: testConn,
			},
			args: args{
				ctx: context.Background(),
				id:  p.ID,
			},
			wantErr: false,
		},
		{
			name: "restore non-existing profile",
			fields: fields{
				postgresDB: testConn,
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
			r := &ProfileRepository{
				postgresDB: tt.fields.postgresDB,
			}
			if err := r.Restore(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Restore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
