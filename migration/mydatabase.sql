

-- DROP table if exists user;
-- DROP table if exists post;
-- DROP table if exists comment;
-- DROP table if exists category;
-- Drop table if exists post_category;


CREATE TABLE  IF NOT EXISTS user (
    user_id integer primary key,
    nickname varchar(50),
    email varchar(50),
    password varchar(256),
    role varchar(15)
);

CREATE TABLE  IF NOT EXISTS post (
    post_id integer primary key,
    user_id integer not null,
    title TEXT,
    content TEXT,
    image TEXT,
    like_count integer not null default 0,
    dislike_count integer  not null default 0,
    created_at text not null default (datetime('now')),
    FOREIGN KEY (user_id) REFERENCES user(user_id)
);

CREATE TABLE  IF NOT EXISTS comment(
 comment_id integer primary key,
 post_id integer not null,
 user_id integer not null,
 content TEXT, 
 like_count integer,
 dislike_count integer,
 FOREIGN KEY (post_id)REFERENCES post(post_id) on delete CASCADE,
 FOREIGN KEY (user_id)REFERENCES user(user_id) on delete CASCADE
);


CREATE TABLE  IF NOT EXISTS category(
   category_id integer primary key, 
   category_name text  
);

CREATE TABLE IF NOT EXISTS post_category(
    category_id integer not null,
    post_id integer not null,
    FOREIGN KEY (category_id) REFERENCES category(category_id),
    FOREIGN KEY (post_id) REFERENCES post(post_id) on delete CASCADE
);

CREATE TABLE IF NOT EXISTS report(
    user_id integer not null,
    post_id integer not null,
    FOREIGN KEY (post_id) REFERENCES post(post_id) on delete CASCADE,
    FOREIGN KEY (user_id)REFERENCES user(user_id) on delete CASCADE
);