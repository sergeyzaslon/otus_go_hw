-- +goose Up
CREATE TABLE events (
    id TEXT NOT NULL,
    title TEXT NOT NULL,
    "date" timestamp NOT NULL,
    duration INT NOT NULL,
    description TEXT,
    user_id TEXT NOT NULL DEFAULT '',
    notify_before INT NOT NULL DEFAULT 0,
    notify_at timestamp NOT NULL,
    notified BOOLEAN DEFAULT false
);

CREATE UNIQUE INDEX events_id ON events(id);

--- INSERT INTO events (id, title, date, duration, user_id, notify_before, notify_at) VALUES 
---    ('fcfa069c-28a9-48d6-b48d-befc8133f2b4', 'Event 1', '2022-01-06 11:00:00', 900, 'U1', 0,   '2022-01-06 11:00:00'),
---    ('87a876ce-4488-45a9-bd36-ca72368a7185', 'Event 2', '2022-01-06 12:00:00', 900, 'U2', 900, '2022-01-06 11:45:00'),
---    ('5753e882-91e0-4e1a-a827-eef8d8271e50', 'Event 3', '2022-01-06 11:45:00', 900, 'U3', 0,   '2022-01-06 11:45:00')
---;

-- +goose Down
DROP TABLE events;