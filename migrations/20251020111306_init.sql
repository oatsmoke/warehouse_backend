-- Create "departments" table
CREATE TABLE "public"."departments" (
  "id" bigserial NOT NULL,
  "title" character varying(100) NOT NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "departments_title_key" UNIQUE ("title")
);
-- Create "employees" table
CREATE TABLE "public"."employees" (
  "id" bigserial NOT NULL,
  "last_name" character varying(100) NOT NULL,
  "first_name" character varying(100) NOT NULL,
  "middle_name" character varying(100) NOT NULL,
  "phone" character varying(100) NOT NULL,
  "department" bigint NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "employees_phone_key" UNIQUE ("phone"),
  CONSTRAINT "employees_department_fkey" FOREIGN KEY ("department") REFERENCES "public"."departments" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT
);
-- Create index "idx_employees_department" to table: "employees"
CREATE INDEX "idx_employees_department" ON "public"."employees" ("department");
-- Create "categories" table
CREATE TABLE "public"."categories" (
  "id" bigserial NOT NULL,
  "title" character varying(100) NOT NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "categories_title_key" UNIQUE ("title")
);
-- Create "profiles" table
CREATE TABLE "public"."profiles" (
  "id" bigserial NOT NULL,
  "title" character varying(100) NOT NULL,
  "category" bigint NOT NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "profiles_title_key" UNIQUE ("title"),
  CONSTRAINT "profiles_category_fkey" FOREIGN KEY ("category") REFERENCES "public"."categories" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT
);
-- Create index "idx_profiles_category" to table: "profiles"
CREATE INDEX "idx_profiles_category" ON "public"."profiles" ("category");
-- Create "equipments" table
CREATE TABLE "public"."equipments" (
  "id" bigserial NOT NULL,
  "serial_number" character varying(100) NOT NULL,
  "profile" bigint NOT NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "equipments_serial_number_key" UNIQUE ("serial_number"),
  CONSTRAINT "equipments_profile_fkey" FOREIGN KEY ("profile") REFERENCES "public"."profiles" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT
);
-- Create index "idx_equipments_profile" to table: "equipments"
CREATE INDEX "idx_equipments_profile" ON "public"."equipments" ("profile");
-- Create "companies" table
CREATE TABLE "public"."companies" (
  "id" bigserial NOT NULL,
  "title" character varying(100) NOT NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "companies_title_key" UNIQUE ("title")
);
-- Create "contracts" table
CREATE TABLE "public"."contracts" (
  "id" bigserial NOT NULL,
  "number" character varying(100) NOT NULL,
  "address" character varying(100) NOT NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "contracts_number_key" UNIQUE ("number")
);
-- Create "locations" table
CREATE TABLE "public"."locations" (
  "id" bigserial NOT NULL,
  "equipment" bigint NOT NULL,
  "employee" bigint NOT NULL,
  "company" bigint NOT NULL,
  "move_at" timestamptz NOT NULL DEFAULT now(),
  "move_code" character varying(100) NOT NULL,
  "move_type" character varying(100) NULL,
  "price" character varying(100) NULL,
  "from_department" bigint NULL,
  "from_employee" bigint NULL,
  "from_contract" bigint NULL,
  "to_department" bigint NULL,
  "to_employee" bigint NULL,
  "to_contract" bigint NULL,
  "comment" character varying(100) NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "locations_company_fkey" FOREIGN KEY ("company") REFERENCES "public"."companies" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT,
  CONSTRAINT "locations_employee_fkey" FOREIGN KEY ("employee") REFERENCES "public"."employees" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT,
  CONSTRAINT "locations_equipment_fkey" FOREIGN KEY ("equipment") REFERENCES "public"."equipments" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT,
  CONSTRAINT "locations_from_contract_fkey" FOREIGN KEY ("from_contract") REFERENCES "public"."contracts" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT,
  CONSTRAINT "locations_from_department_fkey" FOREIGN KEY ("from_department") REFERENCES "public"."departments" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT,
  CONSTRAINT "locations_from_employee_fkey" FOREIGN KEY ("from_employee") REFERENCES "public"."employees" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT,
  CONSTRAINT "locations_to_contract_fkey" FOREIGN KEY ("to_contract") REFERENCES "public"."contracts" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT,
  CONSTRAINT "locations_to_department_fkey" FOREIGN KEY ("to_department") REFERENCES "public"."departments" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT,
  CONSTRAINT "locations_to_employee_fkey" FOREIGN KEY ("to_employee") REFERENCES "public"."employees" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT
);
-- Create index "idx_locations_company" to table: "locations"
CREATE INDEX "idx_locations_company" ON "public"."locations" ("company");
-- Create index "idx_locations_employee" to table: "locations"
CREATE INDEX "idx_locations_employee" ON "public"."locations" ("employee");
-- Create index "idx_locations_equipment" to table: "locations"
CREATE INDEX "idx_locations_equipment" ON "public"."locations" ("equipment");
-- Create index "idx_locations_from_contract" to table: "locations"
CREATE INDEX "idx_locations_from_contract" ON "public"."locations" ("from_contract");
-- Create index "idx_locations_from_department" to table: "locations"
CREATE INDEX "idx_locations_from_department" ON "public"."locations" ("from_department");
-- Create index "idx_locations_from_employee" to table: "locations"
CREATE INDEX "idx_locations_from_employee" ON "public"."locations" ("from_employee");
-- Create index "idx_locations_move_at" to table: "locations"
CREATE INDEX "idx_locations_move_at" ON "public"."locations" ("move_at");
-- Create index "idx_locations_to_contract" to table: "locations"
CREATE INDEX "idx_locations_to_contract" ON "public"."locations" ("to_contract");
-- Create index "idx_locations_to_department" to table: "locations"
CREATE INDEX "idx_locations_to_department" ON "public"."locations" ("to_department");
-- Create index "idx_locations_to_employee" to table: "locations"
CREATE INDEX "idx_locations_to_employee" ON "public"."locations" ("to_employee");
-- Create "replaces" table
CREATE TABLE "public"."replaces" (
  "id" bigserial NOT NULL,
  "move_in" bigint NOT NULL,
  "move_out" bigint NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "replaces_move_in_fkey" FOREIGN KEY ("move_in") REFERENCES "public"."locations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "replaces_move_out_fkey" FOREIGN KEY ("move_out") REFERENCES "public"."locations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_replaces_move_in" to table: "replaces"
CREATE INDEX "idx_replaces_move_in" ON "public"."replaces" ("move_in");
-- Create index "idx_replaces_move_out" to table: "replaces"
CREATE INDEX "idx_replaces_move_out" ON "public"."replaces" ("move_out");
-- Create "users" table
CREATE TABLE "public"."users" (
  "id" bigserial NOT NULL,
  "username" character varying(100) NOT NULL,
  "password_hash" character varying(100) NOT NULL,
  "email" character varying(100) NOT NULL,
  "role" character varying(100) NOT NULL,
  "enabled" boolean NOT NULL DEFAULT true,
  "last_login_at" timestamptz NULL,
  "employee" bigint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "users_username_key" UNIQUE ("username"),
  CONSTRAINT "users_employee_fkey" FOREIGN KEY ("employee") REFERENCES "public"."employees" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT
);
-- Create index "idx_users_employee" to table: "users"
CREATE INDEX "idx_users_employee" ON "public"."users" ("employee");
