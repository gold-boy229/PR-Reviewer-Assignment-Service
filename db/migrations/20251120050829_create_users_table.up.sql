create table users (
	user_id text primary key,
	username text not null,
	team_name text not null,
	is_active bool not null
);