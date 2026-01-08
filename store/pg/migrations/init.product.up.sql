CREATE TABLE PRODUCTS (
    id bigserial primary key,
	product_name  text not null ,
	product_group text,
	image        bytea,
	stock        bigint,
	price        bigint
)