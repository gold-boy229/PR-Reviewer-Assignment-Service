alter table pull_requests__M2M__users 
drop constraint fk_reviewer_id__users;

alter table pull_requests__M2M__users 
drop constraint fk_link_pull_requests;

drop table pull_requests__M2M__users;