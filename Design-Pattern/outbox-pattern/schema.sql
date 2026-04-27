create table orders (
    id string primary key,
    Product string not null,
    Quantity int not null
);

create table outbox (
    id uuid primary key default gen_random_uuid(),
    topic varchar(255) not null,
    message jsonb not null,
    state varchar(20) not null default 'pending',
    created_at timestamptz not null default now(),
    processed_at timestamptz
)