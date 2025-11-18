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

func truncateProfiles(t *testing.T, testDB *pgxpool.Pool) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const query = `
		TRUNCATE profiles, categories
		RESTART IDENTITY CASCADE;`

	if _, err := testDB.Exec(ctx, query); err != nil {
		t.Fatalf("failed to truncate profile: %v", err)
	}
}

func addTestProfile(t *testing.T, testDB *pgxpool.Pool) *model.Profile {
	t.Helper()
	c := addTestCategory(t, testDB)
	p := new(model.Profile)

	const query = `
		INSERT INTO profiles (title, category_id)
		VALUES ($1, $2)
		RETURNING id, title;`

	if err := testDB.QueryRow(t.Context(), query, generate.RandString(10), c.ID).
		Scan(&p.ID, &p.Title); err != nil {
		t.Fatalf("failed to insert test profile: %v", err)
	}

	p.Category = c

	return p
}

func addTestDeletedProfile(t *testing.T, testDB *pgxpool.Pool) *model.Profile {
	t.Helper()
	c := addTestCategory(t, testDB)
	p := new(model.Profile)

	const query = `
		INSERT INTO profiles (title, category_id, deleted_at)
		VALUES ($1, $2, now())
		RETURNING id, title, deleted_at;`

	if err := testDB.QueryRow(t.Context(), query, generate.RandString(10), c.ID).
		Scan(&p.ID, &p.Title, &p.DeletedAt); err != nil {
		t.Fatalf("failed to insert test profile: %v", err)
	}

	p.Category = c

	return p
}

func TestNewProfileRepository(t *testing.T) {
	testDB := postgresql.ConnectTest()
	defer testDB.Close()
	q := queries.New(testDB)

	type args struct {
		queries queries.Querier
	}
	tests := []struct {
		name string
		args args
		want *ProfileRepository
	}{
		{
			name: "create profile repository",
			args: args{
				queries: q,
			},
			want: NewProfileRepository(q),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewProfileRepository(tt.args.queries); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProfileRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProfileRepository_Create(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateProfiles(t, testDB)
		testDB.Close()
	})
	q := queries.New(testDB)
	c := addTestCategory(t, testDB)

	type fields struct {
		queries queries.Querier
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
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				profile: &model.Profile{
					Title: "test profile",
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
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				profile: &model.Profile{
					Title: "test profile",
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
				queries: tt.fields.queries,
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
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateProfiles(t, testDB)
		testDB.Close()
	})
	q := queries.New(testDB)
	p := addTestProfile(t, testDB)

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
		want    *model.Profile
		wantErr bool
	}{
		{
			name: "read profile",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				id:  p.ID,
			},
			want:    p,
			wantErr: false,
		},
		{
			name: "read non-existing profile",
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
			r := &ProfileRepository{
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

func TestProfileRepository_Update(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateProfiles(t, testDB)
		testDB.Close()
	})
	q := queries.New(testDB)
	p := addTestProfile(t, testDB)

	type fields struct {
		queries queries.Querier
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
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				profile: &model.Profile{
					ID:    p.ID,
					Title: "updated profile",
					Category: &model.Category{
						ID: p.Category.ID,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "update non-existing profile",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				profile: &model.Profile{
					ID:    999,
					Title: "update non-existing profile",
					Category: &model.Category{
						ID: p.Category.ID,
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ProfileRepository{
				queries: tt.fields.queries,
			}
			if err := r.Update(tt.args.ctx, tt.args.profile); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProfileRepository_Delete(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateProfiles(t, testDB)
		testDB.Close()
	})
	q := queries.New(testDB)
	p := addTestProfile(t, testDB)

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
			name: "delete profile",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				id:  p.ID,
			},
			wantErr: false,
		},
		{
			name: "delete non-existing profile",
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
			r := &ProfileRepository{
				queries: tt.fields.queries,
			}
			if err := r.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProfileRepository_Restore(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateProfiles(t, testDB)
		testDB.Close()
	})
	q := queries.New(testDB)
	p := addTestProfile(t, testDB)
	dp := addTestDeletedProfile(t, testDB)

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
			name: "restore profile",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				id:  dp.ID,
			},
			wantErr: false,
		},
		{
			name: "restore non-existing profile",
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
			name: "restore not deleted profile",
			fields: fields{
				queries: q,
			},
			args: args{
				ctx: t.Context(),
				id:  p.ID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ProfileRepository{
				queries: tt.fields.queries,
			}
			if err := r.Restore(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Restore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProfileRepository_List(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateProfiles(t, testDB)
		testDB.Close()
	})
	q := queries.New(testDB)
	p := addTestProfile(t, testDB)
	dp := addTestDeletedProfile(t, testDB)

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
		want    []*model.Profile
		want1   int64
		wantErr bool
	}{
		{
			name: "list profiles without deleted",
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
			want:    []*model.Profile{p},
			want1:   1,
			wantErr: false,
		},
		{
			name: "list profiles with deleted",
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
			want:    []*model.Profile{p, dp},
			want1:   2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ProfileRepository{
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
