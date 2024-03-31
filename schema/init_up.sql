create table categories
(
    id    serial       not null primary key,
    title varchar(100) not null unique
);
create table profiles
(
    id       serial                             not null primary key,
    title    varchar(100)                       not null unique,
    category integer references categories (id) not null
);
create table equipments
(
    id            serial                           not null primary key,
    serial_number varchar(100)                     not null unique,
    profile       integer references profiles (id) not null,
    is_deleted    boolean                          not null default false
);
create table departments
(
    id         serial       not null primary key,
    title      varchar(100) not null unique,
    is_deleted boolean      not null default false
);
create table employees
(
    id                 serial                   not null primary key,
    name               varchar(100)             not null,
    phone              varchar(100)             not null unique,
    email              varchar(100)             not null,
    password           varchar(100)             not null,
    hash               varchar(100)             not null,
    registration_date  timestamp with time zone not null,
    authorization_date timestamp with time zone not null,
    activate           boolean                  not null default false,
    hidden             boolean                  not null default false,
    department         integer references departments (id),
    role               varchar(100)             not null default 'USER',
    is_deleted         boolean                  not null default false
);
create table contracts
(
    id         serial       not null primary key,
    number     varchar(100) not null unique,
    address    varchar(100) not null,
    is_deleted boolean      not null default false
);
create table companies
(
    id         serial       not null primary key,
    title      varchar(100) not null unique,
    is_deleted boolean      not null default false
);
create table locations
(
    id              serial                             not null primary key,
    date            timestamp with time zone           not null,
    code            varchar(100)                       not null,
    equipment       integer references equipments (id) not null,
    employee        integer references employees (id)  not null,
    company         integer references companies (id)  not null,
    from_department integer references departments (id),
    from_employee   integer references employees (id),
    from_contract   integer references contracts (id),
    to_department   integer references departments (id),
    to_employee     integer references employees (id),
    to_contract     integer references contracts (id),
    transfer_type   varchar(100),
    price           integer
);
create table replaces
(
    id            serial                                         not null primary key,
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
        '313233343536373840bd001563085fc35165329ea1ff5c5ecbdbbeef',
        '',
        now(),
        now(),
        true,
        true,
        'ADMIN');