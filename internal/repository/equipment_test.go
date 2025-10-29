package repository

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oatsmoke/warehouse_backend/internal/dto"
	"github.com/oatsmoke/warehouse_backend/internal/lib/generate"
	"github.com/oatsmoke/warehouse_backend/internal/lib/postgresql"
	"github.com/oatsmoke/warehouse_backend/internal/model"
)

func truncateEquipments(t *testing.T, testDB *pgxpool.Pool) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const query = `
		TRUNCATE equipments, profiles, categories
		RESTART IDENTITY CASCADE;`

	if _, err := testDB.Exec(ctx, query); err != nil {
		t.Fatalf("failed to truncate equipment: %v", err)
	}
}

func addTestEquipment(t *testing.T, testDB *pgxpool.Pool) *model.Equipment {
	t.Helper()
	p := addTestProfile(t, testDB)
	e := new(model.Equipment)

	const query = `
		INSERT INTO equipments (serial_number, profile)
		VALUES ($1, $2)
		RETURNING id, serial_number;`

	if err := testDB.QueryRow(t.Context(), query, generate.RandString(10), p.ID).
		Scan(&e.ID, &e.SerialNumber); err != nil {
		t.Fatalf("failed to insert test equipment: %v", err)
	}

	e.Profile = p

	return e
}

func addTestDeletedEquipment(t *testing.T, testDB *pgxpool.Pool) *model.Equipment {
	t.Helper()
	p := addTestProfile(t, testDB)
	e := new(model.Equipment)

	const query = `
		INSERT INTO equipments (serial_number, profile, deleted_at)
		VALUES ($1, $2, now())
		RETURNING id, serial_number, deleted_at;`

	if err := testDB.QueryRow(t.Context(), query, generate.RandString(10), p.ID).
		Scan(&e.ID, &e.SerialNumber, &e.DeletedAt); err != nil {
		t.Fatalf("failed to insert test equipment: %v", err)
	}

	e.Profile = p

	return e
}

func TestNewEquipmentRepository(t *testing.T) {
	testDB := postgresql.ConnectTest()
	defer testDB.Close()

	type args struct {
		postgresDB *pgxpool.Pool
	}
	tests := []struct {
		name string
		args args
		want *EquipmentRepository
	}{
		{
			name: "create equipment repository",
			args: args{
				postgresDB: testDB,
			},
			want: NewEquipmentRepository(testDB),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEquipmentRepository(tt.args.postgresDB); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEquipmentRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEquipmentRepository_Create(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateEquipments(t, testDB)
		testDB.Close()
	})
	p := addTestProfile(t, testDB)

	type fields struct {
		postgresDB *pgxpool.Pool
	}
	type args struct {
		ctx       context.Context
		equipment *model.Equipment
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "create equipment",
			fields: fields{
				postgresDB: testDB,
			},
			args: args{
				ctx: t.Context(),
				equipment: &model.Equipment{
					SerialNumber: "test equipment",
					Profile: &model.Profile{
						ID: p.ID,
					},
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "create duplicate equipment",
			fields: fields{
				postgresDB: testDB,
			},
			args: args{
				ctx: t.Context(),
				equipment: &model.Equipment{
					SerialNumber: "test equipment",
					Profile: &model.Profile{
						ID: p.ID,
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &EquipmentRepository{
				postgresDB: tt.fields.postgresDB,
			}
			got, err := r.Create(tt.args.ctx, tt.args.equipment)
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

func TestEquipmentRepository_Read(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateEquipments(t, testDB)
		testDB.Close()
	})
	e := addTestEquipment(t, testDB)

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
		want    *model.Equipment
		wantErr bool
	}{
		{
			name: "read equipment",
			fields: fields{
				postgresDB: testDB,
			},
			args: args{
				ctx: t.Context(),
				id:  e.ID,
			},
			want:    e,
			wantErr: false,
		},
		{
			name: "read non-existing equipment",
			fields: fields{
				postgresDB: testDB,
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
			r := &EquipmentRepository{
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

func TestEquipmentRepository_Update(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateEquipments(t, testDB)
		testDB.Close()
	})
	e := addTestEquipment(t, testDB)

	type fields struct {
		postgresDB *pgxpool.Pool
	}
	type args struct {
		ctx       context.Context
		equipment *model.Equipment
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "update equipment",
			fields: fields{
				postgresDB: testDB,
			},
			args: args{
				ctx: t.Context(),
				equipment: &model.Equipment{
					ID:           e.ID,
					SerialNumber: "updated_equipment",
					Profile: &model.Profile{
						ID: e.Profile.ID,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "update non-existing equipment",
			fields: fields{
				postgresDB: testDB,
			},
			args: args{
				ctx: t.Context(),
				equipment: &model.Equipment{
					ID:           999,
					SerialNumber: "update non-existing equipment",
					Profile: &model.Profile{
						ID: e.Profile.ID,
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &EquipmentRepository{
				postgresDB: tt.fields.postgresDB,
			}
			if err := r.Update(tt.args.ctx, tt.args.equipment); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEquipmentRepository_Delete(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateEquipments(t, testDB)
		testDB.Close()
	})
	e := addTestEquipment(t, testDB)

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
			name: "delete equipment",
			fields: fields{
				postgresDB: testDB,
			},
			args: args{
				ctx: t.Context(),
				id:  e.ID,
			},
			wantErr: false,
		},
		{
			name: "delete non-existing equipment",
			fields: fields{
				postgresDB: testDB,
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
			r := &EquipmentRepository{
				postgresDB: tt.fields.postgresDB,
			}
			if err := r.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEquipmentRepository_Restore(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateEquipments(t, testDB)
		testDB.Close()
	})
	e := addTestEquipment(t, testDB)
	de := addTestDeletedEquipment(t, testDB)

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
			name: "restore equipment",
			fields: fields{
				postgresDB: testDB,
			},
			args: args{
				ctx: t.Context(),
				id:  de.ID,
			},
			wantErr: false,
		},
		{
			name: "restore non-existing equipment",
			fields: fields{
				postgresDB: testDB,
			},
			args: args{
				ctx: t.Context(),
				id:  999,
			},
			wantErr: true,
		},
		{
			name: "restore not deleted equipment",
			fields: fields{
				postgresDB: testDB,
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
			r := &EquipmentRepository{
				postgresDB: tt.fields.postgresDB,
			}
			if err := r.Restore(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Restore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEquipmentRepository_List(t *testing.T) {
	testDB := postgresql.ConnectTest()
	t.Cleanup(func() {
		truncateEquipments(t, testDB)
		testDB.Close()
	})
	e := addTestEquipment(t, testDB)
	de := addTestDeletedEquipment(t, testDB)

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
		want    []*model.Equipment
		want1   int
		wantErr bool
	}{
		{
			name: "list equipments without deleted",
			fields: fields{
				postgresDB: testDB,
			},
			args: args{
				ctx: t.Context(),
				qp:  &dto.QueryParams{},
			},
			want:    []*model.Equipment{e},
			want1:   1,
			wantErr: false,
		},
		{
			name: "list equipments with deleted",
			fields: fields{
				postgresDB: testDB,
			},
			args: args{
				ctx: t.Context(),
				qp: &dto.QueryParams{
					WithDeleted: true,
				},
			},
			want:    []*model.Equipment{e, de},
			want1:   2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &EquipmentRepository{
				postgresDB: tt.fields.postgresDB,
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
