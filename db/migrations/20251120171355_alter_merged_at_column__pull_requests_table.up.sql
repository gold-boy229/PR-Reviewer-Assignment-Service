alter table pull_requests
alter column merged_at drop default;

comment on column pull_requests.status is 'Enum: [OPEN, MERGED]';