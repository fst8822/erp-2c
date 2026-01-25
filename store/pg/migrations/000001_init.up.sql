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

CREATE TABLE DELIVERY (
                        id bigserial primary key,
                        recipient_goods text NOT NULL,
                        address         text not null,
                        status_delivery text NOT NULL,
                        created_At      timestamp NOT NULL,
                        total_amount bigint not null
);

CREATE TABLE DELIVERY_PRODUCT (
                        id           bigserial  primary key,
                        delivery_id  bigint  not null ,
                        product_id   bigint  not null ,
                        quantity     bigint not null,
                        unit_price   bigint not null,
                        total_amount bigint not null,
                        FOREIGN KEY (delivery_id) REFERENCES DELIVERY(id),
                        FOREIGN KEY (product_id) REFERENCES PRODUCTS(id)
)

