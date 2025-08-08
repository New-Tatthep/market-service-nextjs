CREATE TABLE IF NOT EXISTS employee (
    user_code VARCHAR PRIMARY KEY,
    user_name VARCHAR NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(150) NOT NULL,
    status VARCHAR(10) NOT NULL,
    profile_image JSONB NULL,
    update_code VARCHAR(100) NOT NULL,
    update_time INT8 NOT NULL,
    mobile_no VARCHAR(20) NOT NULL,
    password VARCHAR(255) NULL
);

