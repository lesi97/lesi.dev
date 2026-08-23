create schema if not exists logs;

create table if not exists logs.date_meme_clicks (
    route text not null,
    ip text not null,
    click_date date not null default ((now() at time zone 'utc')::date),
    yes_clicks integer not null default 0 check (yes_clicks >= 0),
    no_clicks integer not null default 0 check (no_clicks >= 0),
    secret_endings integer not null default 0 check (secret_endings >= 0),
    user_agent text not null default '',
    first_clicked_at timestamptz not null default now(),
    last_clicked_at timestamptz not null default now(),

    constraint date_meme_clicks_pkey primary key (route, ip, click_date),
    constraint date_meme_clicks_action_count_check check (yes_clicks > 0 or no_clicks > 0)
);

create index if not exists date_meme_clicks_route_click_date_idx
    on logs.date_meme_clicks (route, click_date desc);

create index if not exists date_meme_clicks_last_clicked_at_idx
    on logs.date_meme_clicks (last_clicked_at desc);
