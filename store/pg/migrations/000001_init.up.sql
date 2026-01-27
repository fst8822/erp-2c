CREATE TABLE PRODUCTS (
                          id bigserial primary key,
                          product_name  text not null ,
                          product_group text,
                          image        bytea,
                          stock        bigint not null DEFAULT 1,
                          price        bigint not null DEFAULT 1
);
CREATE TABLE USERS (
                       id        bigserial  primary key,
                       first_name text not null,
                       email      text not null unique,
                       login      text not null unique,
                       password   text not null,
                       user_role  text not null
);

CREATE TABLE "delivery" (
                        id bigserial primary key,
                        recipient       text NOT NULL,
                        address         text not null,
                        status          text NOT NULL,
                        created_At      timestamp NOT NULL
);

CREATE TABLE delivery_items (
    id bigserial primary key,
    Delivery_id bigint NOT NULL REFERENCES DELIVERY(id),
    Product_id bigint NOT NULL REFERENCES products(id),
    item_price bigint NOT NULL,
    quantity bigint NOT NULL
)

