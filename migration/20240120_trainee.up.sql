create table users (
    uuid uuid DEFAULT uuid_generate_v4() primary key,
    username varchar(255) not null unique,
    password varchar(255) not null,
    created_at timestamp not null default now()
);

create table wallets (
    uuid uuid DEFAULT uuid_generate_v4() primary key,
    balance float8 not null
);

create table transactions (
    uuid uuid default uuid_generate_v4() primary key,
    sender uuid not null ,
    receiver uuid not null,
    created_at timestamp not null default now(),
    amount float8 not null,
    foreign key (sender) references wallets(uuid) on delete cascade,
    foreign key (receiver) references wallets(uuid) on delete cascade
);