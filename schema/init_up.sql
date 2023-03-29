CREATE TABLE categories
(
    id    serial       not null primary key,
    title varchar(100) not null unique
);
CREATE TABLE profiles
(
    id       serial                             not null primary key,
    title    varchar(100)                       not null unique,
    category integer references categories (id) not null
);
CREATE TABLE equipments
(
    id            serial                           not null primary key,
    serial_number varchar(100)                     not null unique,
    profile       integer references profiles (id) not null,
    is_deleted    boolean                          not null default false
);
CREATE TABLE departments
(
    id         serial       not null primary key,
    title      varchar(100) not null unique,
    is_deleted boolean      not null default false
);
CREATE TABLE employees
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
CREATE TABLE contracts
(
    id         serial       not null primary key,
    number     varchar(100) not null unique,
    address    varchar(100) not null,
    is_deleted boolean      not null default false
);
CREATE TABLE locations
(
    id            serial                             not null primary key,
    date          timestamp with time zone           not null,
    code          varchar(100)                       not null,
    equipment     integer references equipments (id) not null,
    employee      integer references employees (id)  not null,
    to_department integer references departments (id),
    to_employee   integer references employees (id),
    to_contract   integer references contracts (id)
);
INSERT INTO employees (name,
                       phone,
                       email,
                       password,
                       hash,
                       registration_date,
                       authorization_date,
                       activate,
                       hidden,
                       role)
VALUES ('Администратор',
        'root',
        'root@root.ru',
        '313233343536373840bd001563085fc35165329ea1ff5c5ecbdbbeef',
        '',
        now(),
        now(),
        true,
        true,
        'ADMIN');