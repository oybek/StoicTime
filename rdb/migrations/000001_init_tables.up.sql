
CREATE TABLE act (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  name VARCHAR(100) NOT NULL,

  UNIQUE (user_id, name)
);

CREATE TABLE act_log (
  message_id BIGINT NOT NULL,
  user_id BIGINT NOT NULL,
  name VARCHAR(100) NOT NULL,
  start_time TIMESTAMP NOT NULL,
  end_time TIMESTAMP NOT NULL
);

CREATE UNIQUE INDEX one_zero_endtime_per_user
  ON act_log (user_id)
  WHERE end_time = '0001-01-01Z00:00:00';
