ALTER TABLE tasks DROP COLUMN deadline;

ALTER TABLE tasks DROP COLUMN board_id;

ALTER TABLE tasks
    ADD COLUMN status TEXT;