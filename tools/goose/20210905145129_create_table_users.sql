-- +goose Up
-- +goose StatementBegin
CREATE TABLE users(
  id SERIAL,
  firstname VARCHAR(255),
  lastname VARCHAR(255),
  birthday VARCHAR(255),
  active BOOLEAN
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP table users;
-- +goose StatementEnd
