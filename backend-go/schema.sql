CREATE EXTENSION IF NOT EXISTS "pg_catalog.plpgsql";

DROP TABLE IF EXISTS message_responses CASCADE;
DROP TABLE IF EXISTS message_deliveries CASCADE;
DROP TABLE IF EXISTS invitations CASCADE;
DROP TABLE IF EXISTS messages CASCADE;
DROP TABLE IF EXISTS students CASCADE;
DROP TABLE IF EXISTS classrooms CASCADE;
DROP TABLE IF EXISTS jwt_denylists CASCADE;
DROP TABLE IF EXISTS posts CASCADE;
DROP TABLE IF EXISTS teachers CASCADE;

CREATE TABLE teachers (
  id BIGSERIAL PRIMARY KEY,
  email VARCHAR NOT NULL DEFAULT '',
  encrypted_password VARCHAR NOT NULL DEFAULT '',
  reset_password_token VARCHAR,
  reset_password_sent_at TIMESTAMP,
  remember_created_at TIMESTAMP,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  name VARCHAR
);

CREATE UNIQUE INDEX index_teachers_on_email
  ON teachers (email);

CREATE UNIQUE INDEX index_teachers_on_reset_password_token
  ON teachers (reset_password_token);

CREATE TABLE classrooms (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR,
  teacher_id BIGINT NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

CREATE INDEX index_classrooms_on_teacher_id
  ON classrooms (teacher_id);

CREATE TABLE invitations (
  id BIGSERIAL PRIMARY KEY,
  token VARCHAR NOT NULL,
  email VARCHAR NOT NULL,
  classroom_id BIGINT NOT NULL,
  used BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  expires_at TIMESTAMP,
  used_at TIMESTAMP
);

CREATE INDEX index_invitations_on_classroom_id
  ON invitations (classroom_id);

CREATE UNIQUE INDEX index_invitations_on_token
  ON invitations (token);

CREATE TABLE jwt_denylists (
  id BIGSERIAL PRIMARY KEY,
  jti VARCHAR,
  exp TIMESTAMP,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

CREATE INDEX index_jwt_denylists_on_jti
  ON jwt_denylists (jti);

CREATE TABLE posts (
  id BIGSERIAL PRIMARY KEY,
  title VARCHAR,
  body TEXT,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

CREATE TABLE messages (
  id BIGSERIAL PRIMARY KEY,
  title VARCHAR,
  content TEXT,
  status INTEGER NOT NULL DEFAULT 0,
  published_at TIMESTAMP,
  deadline DATE,
  classroom_id BIGINT NOT NULL,
  teacher_id BIGINT NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  target_all BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE INDEX index_messages_on_classroom_id_and_status
  ON messages (classroom_id, status);

CREATE INDEX index_messages_on_classroom_id
  ON messages (classroom_id);

CREATE INDEX index_messages_on_deadline
  ON messages (deadline);

CREATE INDEX index_messages_on_published_at
  ON messages (published_at);

CREATE INDEX index_messages_on_teacher_id
  ON messages (teacher_id);

CREATE TABLE students (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR,
  email VARCHAR,
  classroom_id BIGINT NOT NULL,
  encrypted_password VARCHAR,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

CREATE INDEX index_students_on_classroom_id
  ON students (classroom_id);

CREATE TABLE message_deliveries (
  id BIGSERIAL PRIMARY KEY,
  message_id BIGINT NOT NULL,
  student_id BIGINT NOT NULL,
  confirmed_at TIMESTAMP,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

CREATE INDEX index_message_deliveries_on_confirmed_at
  ON message_deliveries (confirmed_at);

CREATE UNIQUE INDEX index_message_deliveries_on_message_id_and_student_id
  ON message_deliveries (message_id, student_id);

CREATE INDEX index_message_deliveries_on_message_id
  ON message_deliveries (message_id);

CREATE INDEX index_message_deliveries_on_student_id
  ON message_deliveries (student_id);

CREATE TABLE message_responses (
  id BIGSERIAL PRIMARY KEY,
  message_delivery_id BIGINT NOT NULL,
  status INTEGER NOT NULL DEFAULT 0,
  form_data JSONB NOT NULL DEFAULT '{}'::jsonb,
  responded_at TIMESTAMP,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

CREATE INDEX index_message_responses_on_message_delivery_id
  ON message_responses (message_delivery_id);

CREATE INDEX index_message_responses_on_responded_at
  ON message_responses (responded_at);

CREATE INDEX index_message_responses_on_status
  ON message_responses (status);

ALTER TABLE classrooms
  ADD CONSTRAINT fk_classrooms_teacher
  FOREIGN KEY (teacher_id) REFERENCES teachers(id);

ALTER TABLE invitations
  ADD CONSTRAINT fk_invitations_classroom
  FOREIGN KEY (classroom_id) REFERENCES classrooms(id);

ALTER TABLE messages
  ADD CONSTRAINT fk_messages_classroom
  FOREIGN KEY (classroom_id) REFERENCES classrooms(id);

ALTER TABLE messages
  ADD CONSTRAINT fk_messages_teacher
  FOREIGN KEY (teacher_id) REFERENCES teachers(id);

ALTER TABLE students
  ADD CONSTRAINT fk_students_classroom
  FOREIGN KEY (classroom_id) REFERENCES classrooms(id);

ALTER TABLE message_deliveries
  ADD CONSTRAINT fk_message_deliveries_message
  FOREIGN KEY (message_id) REFERENCES messages(id);

ALTER TABLE message_deliveries
  ADD CONSTRAINT fk_message_deliveries_student
  FOREIGN KEY (student_id) REFERENCES students(id);

ALTER TABLE message_responses
  ADD CONSTRAINT fk_message_responses_message_delivery
  FOREIGN KEY (message_delivery_id) REFERENCES message_deliveries(id);