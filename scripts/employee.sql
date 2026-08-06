CREATE
OR REPLACE PROCEDURE initial_employee()
LANGUAGE plpgsql
AS
$$

BEGIN 
 -- 1. ตาราง employee
 IF EXISTS(SELECT FROM pg_catalog.pg_tables WHERE schemaname = 'public' AND tablename = 'employee') THEN
        RAISE NOTICE 'Table public.employee already exists.';
    ELSE
        CREATE TABLE employee (
            user_code VARCHAR NOT NULL,
            user_name VARCHAR NOT NULL,
            user_type varchar(10) NOT NULL,
            password VARCHAR(255) NOT NULL,
            first_name VARCHAR(100) NOT NULL,
            last_name VARCHAR(100) NOT NULL,
            email VARCHAR(150) NOT NULL,
            birth_date          bigint  NULL,
            profile_image      JSONB NULL DEFAULT '{}'::JSONB,
            mobile_no VARCHAR(20)  NULL,
            create_code             varchar(50) NOT NULL,   
            create_time             bigint NOT NULL,
            update_code             varchar(50) NOT NULL,
            update_time             bigint NOT NULL,
            status                  varchar(20) NOT NULL,
            count_login_fail int4 NULL DEFAULT 0,

            CONSTRAINT pk_employee PRIMARY KEY (user_code)
        );   
END IF;

-- 2. ตาราง employee_profile_edit_history
IF EXISTS (SELECT FROM pg_catalog.pg_tables WHERE schemaname = 'public' AND tablename = 'employee_profile_edit_history') THEN
        RAISE NOTICE 'Table public.employee_profile_edit_history already exists.';
    ELSE
        CREATE TABLE IF NOT EXISTS  employee_profile_edit_history (
        history_code  varchar(50) NOT NULL,
        edit_data     jsonb        NULL,
        update_time   bigint       NOT NULL,
        update_code   varchar(100) NOT NULL,
        user_code     varchar(100) NOT NULL,
        CONSTRAINT pk_employee_profile_edit_history PRIMARY KEY (history_code)
    );

    ALTER TABLE employee_profile_edit_history ADD COLUMN IF NOT EXISTS user_code varchar(100) NOT NULL DEFAULT '';

    CREATE INDEX IF NOT EXISTS idx_employee_profile_edit_history_history_code ON employee_profile_edit_history (history_code);
    CREATE INDEX IF NOT EXISTS idx_employee_profile_edit_history_update_code ON employee_profile_edit_history (update_code);
    CREATE INDEX IF NOT EXISTS idx_employee_profile_edit_history_update_time ON employee_profile_edit_history (update_time);   

END IF;

-- 3. ตาราง category (หมวดหมู่สินค้า)
IF EXISTS (SELECT FROM pg_catalog.pg_tables WHERE schemaname = 'public' AND tablename = 'category') THEN
        RAISE NOTICE 'Table public.category already exists.';
    ELSE
        CREATE TABLE category (
            code VARCHAR(50) PRIMARY KEY,
            name VARCHAR(100) NOT NULL,
            description TEXT NULL,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        );
END IF;

-- 4. ตาราง product (สินค้า)
IF EXISTS (SELECT FROM pg_catalog.pg_tables WHERE schemaname = 'public' AND tablename = 'product') THEN
        RAISE NOTICE 'Table public.product already exists.';
    ELSE
        CREATE TABLE product (
            code VARCHAR(50) PRIMARY KEY,
            category_code VARCHAR(50) REFERENCES category(code) ON DELETE SET NULL,
            name VARCHAR(255) NOT NULL,
            description TEXT NULL,
            price DECIMAL(10, 2) NOT NULL,
            status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
            quantity INT8 NOT NULL DEFAULT 0,
            image JSONB NULL,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        );
        
        CREATE INDEX IF NOT EXISTS idx_product_category ON product(category_code);
END IF;

-- 5. ตาราง cart_item (ตะกร้าสินค้า)
IF EXISTS (SELECT FROM pg_catalog.pg_tables WHERE schemaname = 'public' AND tablename = 'cart_item') THEN
        RAISE NOTICE 'Table public.cart_item already exists.';
    ELSE
        CREATE TABLE cart_item (
            id SERIAL PRIMARY KEY,
            user_id VARCHAR(100) NOT NULL,
            product_code VARCHAR(50) REFERENCES product(code) ON DELETE CASCADE,
            quantity INT NOT NULL DEFAULT 1,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            CONSTRAINT unique_user_product UNIQUE (user_id, product_code)
        );

        CREATE INDEX IF NOT EXISTS idx_cart_user ON cart_item(user_id);
END IF;

END
$$;

CALL initial_employee();
DROP procedure if exists initial_employee();