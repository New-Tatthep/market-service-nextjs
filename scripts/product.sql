CREATE TABLE IF NOT EXISTS product (
    code VARCHAR PRIMARY KEY,
    name VARCHAR NOT NULL,
    description VARCHAR NULL,
    price DECIMAL(10, 2) NOT NULL,
    status varchar(20) NOT NULL,
    quantity int8 NOT NULL,
    image jsonb NULL
);