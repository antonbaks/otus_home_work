-- +goose Up
alter table events
    add is_notified smallint default 0 not null;

-- +goose Down
alter table events
    drop column is_notified;