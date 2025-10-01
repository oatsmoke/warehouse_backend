-- Create "departments" table
CREATE TABLE "public"."departments" (
  "id" bigserial NOT NULL,
  "title" character varying(100) NOT NULL,
  "deleted" boolean NOT NULL DEFAULT false,
  PRIMARY KEY ("id"),
  CONSTRAINT "departments_title_key" UNIQUE ("title")
);
-- Create "employees" table
CREATE TABLE "public"."employees" (
  "id" bigserial NOT NULL,
  "name" character varying(100) NOT NULL,
  "phone" character varying(100) NOT NULL,
  "email" character varying(100) NOT NULL,
  "password" character varying(100) NOT NULL,
  "hash" character varying(100) NOT NULL,
  "registration_date" timestamptz NULL,
  "authorization_date" timestamptz NULL,
  "activate" boolean NOT NULL DEFAULT false,
  "hidden" boolean NOT NULL DEFAULT false,
  "department" bigint NULL,
  "role" character varying(100) NOT NULL DEFAULT 'USER',
  "deleted" boolean NOT NULL DEFAULT false,
  PRIMARY KEY ("id"),
  CONSTRAINT "employees_phone_key" UNIQUE ("phone"),
  CONSTRAINT "employees_department_fkey" FOREIGN KEY ("department") REFERENCES "public"."departments" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT
);
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
  "deleted" boolean NOT NULL DEFAULT false,
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
  "deleted" boolean NOT NULL DEFAULT false,
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
  "deleted" boolean NOT NULL DEFAULT false,
  PRIMARY KEY ("id"),
  CONSTRAINT "companies_title_key" UNIQUE ("title")
);
-- Create "contracts" table
CREATE TABLE "public"."contracts" (
  "id" bigserial NOT NULL,
  "number" character varying(100) NOT NULL,
  "address" character varying(100) NOT NULL,
  "deleted" boolean NOT NULL DEFAULT false,
  PRIMARY KEY ("id"),
  CONSTRAINT "contracts_number_key" UNIQUE ("number")
);
-- Create "locations" table
CREATE TABLE "public"."locations" (
  "id" bigserial NOT NULL,
  "date" timestamptz NOT NULL DEFAULT now(),
  "code" character varying(100) NOT NULL,
  "equipment" bigint NOT NULL,
  "employee" bigint NOT NULL,
  "company" bigint NOT NULL,
  "from_department" bigint NULL,
  "from_employee" bigint NULL,
  "from_contract" bigint NULL,
  "to_department" bigint NULL,
  "to_employee" bigint NULL,
  "to_contract" bigint NULL,
  "transfer_type" character varying(100) NULL,
  "price" character varying(100) NULL,
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
-- Create "replaces" table
CREATE TABLE "public"."replaces" (
  "id" bigserial NOT NULL,
  "transfer_from" bigint NOT NULL,
  "transfer_to" bigint NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "replaces_transfer_from_fkey" FOREIGN KEY ("transfer_from") REFERENCES "public"."locations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "replaces_transfer_to_fkey" FOREIGN KEY ("transfer_to") REFERENCES "public"."locations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
