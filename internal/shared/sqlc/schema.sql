CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE jobs (
                      id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
                      name TEXT NOT NULL,
                      handler TEXT NOT NULL,
                      interval INT NOT NULL,
                      status TEXT NOT NULL,
                      last_triggered TIMESTAMP
);
CREATE INDEX idx_jobs_status ON jobs (status);
CREATE INDEX idx_jobs_last_triggered ON jobs (last_triggered);


CREATE TABLE messages (
                          id UUID PRIMARY KEY,
                          content TEXT NOT NULL,--content lenght is importante so make it 200
                          recipient_phone_number TEXT NOT NULL,
                          status TEXT NOT NULL,
                          message_received_id UUID
);