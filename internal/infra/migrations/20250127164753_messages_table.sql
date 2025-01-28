-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE messages
(
    id                     UUID PRIMARY KEY,
    content                VARCHAR(200) NOT NULL,
    recipient_phone_number TEXT         NOT NULL,
    status                 TEXT         NOT NULL,
    message_received_id    UUID
);

INSERT INTO messages (id ,content, recipient_phone_number, status)
VALUES (uuid_generate_v4(),'Lorem ipsum dolor sit amet, consectetur adipiscing elit.', '+901234567890', 'pending'),
       (uuid_generate_v4(),'Lorem ipsum dolor sit amet, consectetur adipiscing elit.', '+905551234567', 'pending'),
       (uuid_generate_v4(),'Lorem ipsum dolor sit amet, consectetur adipiscing elit.', '+902165432187', 'pending'),
       (uuid_generate_v4(),'Lorem ipsum dolor sit amet, consectetur adipiscing elit.', '+903334445566', 'pending'),
       (uuid_generate_v4(),'Lorem ipsum dolor sit amet, consectetur adipiscing elit.', '+904567891234', 'pending'),
       (uuid_generate_v4(),'Lorem ipsum dolor sit amet, consectetur adipiscing elit.', '+905678912345', 'pending'),
       (uuid_generate_v4(),'Lorem ipsum dolor sit amet, consectetur adipiscing elit.', '+906789123456', 'pending'),
       (uuid_generate_v4(),'Lorem ipsum dolor sit amet, consectetur adipiscing elit.', '+907891234567', 'pending'),
       (uuid_generate_v4(),'Lorem ipsum dolor sit amet, consectetur adipiscing elit.', '+908912345678', 'pending'),
       (uuid_generate_v4(),'Lorem ipsum dolor sit amet, consectetur adipiscing elit.', '+909876543210', 'pending');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE messages;
-- +goose StatementEnd

