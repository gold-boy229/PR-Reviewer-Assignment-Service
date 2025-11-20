create table pull_requests__M2M__users (
	pull_request_id text not null,
	reviewer_id text not null
);

alter table pull_requests__M2M__users 
add constraint pk_pr_id__reviewer_id primary key(pull_request_id, reviewer_id);

alter table pull_requests__M2M__users 
add constraint fk_link_pull_requests foreign key (pull_request_id) references pull_requests(pull_request_id);

alter table pull_requests__M2M__users 
add constraint fk_reviewer_id__users foreign key (reviewer_id ) references users(user_id);
