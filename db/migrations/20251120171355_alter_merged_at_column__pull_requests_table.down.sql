alter table pull_requests
alter column merged_at set default now();

comment on column pull_requests.status is null;