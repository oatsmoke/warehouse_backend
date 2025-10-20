create table categories
(
    id         bigserial primary key,
    title      varchar(100) not null unique,
    deleted_at timestamp with time zone
);

create table profiles
(
    id         bigserial primary key,
    title      varchar(100)                                         not null unique,
    category   bigint references categories (id) on delete restrict not null,
    deleted_at timestamp with time zone
);
create index idx_profiles_category on profiles (category);

create table equipments
(
    id            bigserial primary key,
    serial_number varchar(100)                                       not null unique,
    profile       bigint references profiles (id) on delete restrict not null,
    deleted_at    timestamp with time zone
);
create index idx_equipments_profile on equipments (profile);

create table departments
(
    id         bigserial primary key,
    title      varchar(100) not null unique,
    deleted_at timestamp with time zone
);

create table employees
(
    id          bigserial primary key,
    last_name   varchar(100) not null,
    first_name  varchar(100) not null,
    middle_name varchar(100) not null,
    phone       varchar(100) not null unique,
    department  bigint references departments (id) on delete restrict,
    deleted_at  timestamp with time zone
);
create index idx_employees_department on employees (department);

create table users
(
    id            bigserial primary key,
    username      varchar(100) not null unique,
    password_hash varchar(100) not null,
    email         varchar(100) not null,
    role          varchar(100) not null,
    enabled       boolean      not null default true,
    last_login_at timestamp with time zone,
    employee      bigint references employees (id) on delete restrict
);
create index idx_users_employee on users (employee);

create table contracts
(
    id         bigserial primary key,
    number     varchar(100) not null unique,
    address    varchar(100) not null,
    deleted_at timestamp with time zone
);

create table companies
(
    id         bigserial primary key,
    title      varchar(100) not null unique,
    deleted_at timestamp with time zone
);

create table locations
(
    id              bigserial primary key,
    equipment       bigint references equipments (id) on delete restrict not null,
    employee        bigint references employees (id) on delete restrict  not null,
    company         bigint references companies (id) on delete restrict  not null,
    move_at         timestamp with time zone                             not null default now(),
    move_code       varchar(100)                                         not null,
    move_type       varchar(100),
    price           varchar(100),
    from_department bigint references departments (id) on delete restrict,
    from_employee   bigint references employees (id) on delete restrict,
    from_contract   bigint references contracts (id) on delete restrict,
    to_department   bigint references departments (id) on delete restrict,
    to_employee     bigint references employees (id) on delete restrict,
    to_contract     bigint references contracts (id) on delete restrict,
    comment         varchar(100)
);
create index idx_locations_equipment on locations (equipment);
create index idx_locations_employee on locations (employee);
create index idx_locations_company on locations (company);
create index idx_locations_move_at on locations (move_at);
create index idx_locations_from_department on locations (from_department);
create index idx_locations_from_employee on locations (from_employee);
create index idx_locations_from_contract on locations (from_contract);
create index idx_locations_to_department on locations (to_department);
create index idx_locations_to_employee on locations (to_employee);
create index idx_locations_to_contract on locations (to_contract);

create table replaces
(
    id       bigserial primary key,
    move_in  bigint references locations on delete cascade not null,
    move_out bigint references locations on delete cascade not null
);
create index idx_replaces_move_in on replaces (move_in);
create index idx_replaces_move_out on replaces (move_out);