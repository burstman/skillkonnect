-- +goose Up
create table if not exists users(
	id integer primary key,
	email text unique not null,
	password_hash text not null,
	first_name text not null,
	last_name text not null,
	role text not null,
	suspended boolean not null default false,
	approved BOOLEAN DEFAULT FALSE;
	bio TEXT;
	rating REAL DEFAULT 0;
	email_verified_at datetime,
	created_at datetime not null,
	updated_at datetime not null,
	deleted_at datetime
);

-- +goose Down
drop table if exists users;
