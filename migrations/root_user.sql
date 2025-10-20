insert into users (username, password_hash, email, role)
values ('root',
        '$2a$10$sYMtJhDQzFKHk6169kJ4ru8t0phSYEF6NTKjhS9vEewtnXTVcdoIi',
        'root@root.ru',
        'admin')
on conflict (username) do nothing;