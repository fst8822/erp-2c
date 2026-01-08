CREATE TABLE USERS (
                       id        bigserial  primary key,
                       first_name text not null,
                       email     text not null unique,
                       login     text not null unique,
                       password  text not null,
                       user_role  text not null
)