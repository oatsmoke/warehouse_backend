create table categories
(
    id         bigserial primary key,
    title      varchar(100) not null unique,
    deleted_at timestamp with time zone
);

create table profiles
(
    id         bigserial primary key,
    title      varchar(100) not null unique,
    category   bigint       not null references categories (id) on delete restrict,
    deleted_at timestamp with time zone
);

create table equipments
(
    id            bigserial primary key,
    serial_number varchar(100) not null unique,
    profile       bigint       not null references profiles (id) on delete restrict,
    deleted_at    timestamp with time zone
);

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
    date            timestamp with time zone not null default now(),
    code            varchar(100)             not null,
    equipment       bigint                   not null references equipments (id) on delete restrict,
    employee        bigint                   not null references employees (id) on delete restrict,
    company         bigint                   not null references companies (id) on delete restrict,
    from_department bigint references departments (id) on delete restrict,
    from_employee   bigint references employees (id) on delete restrict,
    from_contract   bigint references contracts (id) on delete restrict,
    to_department   bigint references departments (id) on delete restrict,
    to_employee     bigint references employees (id) on delete restrict,
    to_contract     bigint references contracts (id) on delete restrict,
    transfer_type   varchar(100),
    price           varchar(100)
);

create table replaces
(
    id            bigserial primary key,
    transfer_from bigint not null references locations on delete cascade,
    transfer_to   bigint not null references locations on delete cascade
);

create index idx_profiles_category on profiles (category);
create index idx_equipments_profile on equipments (profile);
create index idx_employees_department on employees (department);
create index idx_users_employee on users (employee);
create index idx_locations_equipment on locations (equipment);
create index idx_locations_employee on locations (employee);
create index idx_locations_company on locations (company);

-- insert into employees (name,
--                        phone,
--                        email,
--                        password,
--                        hash,
--                        registration_date,
--                        activate,
--                        hidden,
--                        role)
-- values ('Администратор',
--         'root',
--         'root@root.ru',
--         '$2a$10$sYMtJhDQzFKHk6169kJ4ru8t0phSYEF6NTKjhS9vEewtnXTVcdoIi',
--         '',
--         now(),
--         true,
--         true,
--         'ADMIN');