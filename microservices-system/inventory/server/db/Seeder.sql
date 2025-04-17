
DO $$
    BEGIN
        IF EXISTS (SELECT 1 FROM products) THEN
            RAISE NOTICE 'The table is not empty';
        ELSE
            RAISE NOTICE 'The table is empty';
            INSERT INTO products (name, description, price, stock) VALUES ('Product 1', 'Description 1', 100.00, 10);
            INSERT INTO products (name, description, price, stock) VALUES ('Product 2', 'Description 2', 200.00, 20);
            INSERT INTO products (name, description, price, stock) VALUES ('Product 3', 'Description 3', 300.00, 30);
            insert into users (name, email, password, phone, address, age) values ('John Doe', 'abc1@abc.com', '123456', '1234567890', '123 Main St, New York, NY', 30);
        END IF;
END $$;