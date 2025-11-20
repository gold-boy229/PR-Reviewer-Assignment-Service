create table pull_requests (
	pull_request_id text primary key,
	pull_request_name text not null,
	author_id text not null,
	status text not null,
	created_at timestamptz default now() null,
	merged_at timestamptz default now() null 
);

alter table pull_requests
add constraint fk_author_id__users foreign key (author_id) references users(user_id);