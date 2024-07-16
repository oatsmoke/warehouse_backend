create table categories
(
    id         bigserial not null primary key,
    title      varchar   not null unique,
    is_deleted boolean   not null default false
);
create table profiles
(
    id         bigserial                          not null primary key,
    title      varchar                            not null unique,
    category   integer references categories (id) not null,
    is_deleted boolean                            not null default false
);
create table equipments
(
    id            bigserial                        not null primary key,
    serial_number varchar                          not null unique,
    profile       integer references profiles (id) not null,
    is_deleted    boolean                          not null default false
);
create table departments
(
    id         bigserial not null primary key,
    title      varchar   not null unique,
    is_deleted boolean   not null default false
);
create table employees
(
    id                 bigserial                not null primary key,
    name               varchar                  not null,
    phone              varchar                  not null unique,
    email              varchar                  not null,
    password           varchar                  not null,
    hash               varchar                  not null,
    registration_date  timestamp with time zone not null,
    authorization_date timestamp with time zone not null,
    activate           boolean                  not null default false,
    hidden             boolean                  not null default false,
    department         integer references departments (id),
    role               varchar                  not null default 'USER',
    is_deleted         boolean                  not null default false
);
create table contracts
(
    id         bigserial not null primary key,
    number     varchar   not null unique,
    address    varchar   not null,
    is_deleted boolean   not null default false
);
create table companies
(
    id         bigserial not null primary key,
    title      varchar   not null unique,
    is_deleted boolean   not null default false
);
create table locations
(
    id              bigserial                          not null primary key,
    date            timestamp with time zone           not null,
    code            varchar                            not null,
    equipment       integer references equipments (id) not null,
    employee        integer references employees (id)  not null,
    company         integer references companies (id)  not null,
    from_department integer references departments (id),
    from_employee   integer references employees (id),
    from_contract   integer references contracts (id),
    to_department   integer references departments (id),
    to_employee     integer references employees (id),
    to_contract     integer references contracts (id),
    transfer_type   varchar,
    price           integer
);
create table replaces
(
    id            bigserial                                      not null primary key,
    transfer_from integer references locations on delete cascade not null,
    transfer_to   integer references locations on delete cascade not null
);
insert into employees (name,
                       phone,
                       email,
                       password,
                       hash,
                       registration_date,
                       authorization_date,
                       activate,
                       hidden,
                       role)
values ('Администратор',
        'root',
        'root@root.ru',
        '$2a$10$sYMtJhDQzFKHk6169kJ4ru8t0phSYEF6NTKjhS9vEewtnXTVcdoIi',
        '',
        now(),
        now(),
        true,
        true,
        'ADMIN');