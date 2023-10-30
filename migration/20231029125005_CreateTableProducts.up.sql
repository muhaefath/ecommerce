CREATE TABLE IF NOT EXISTS products (
  id serial PRIMARY KEY,
  user_id bigint NOT NULL,
  sku varchar(255) NOT NULL,
  title varchar(255) NOT NULL,
  description varchar(255) ,
  category varchar(255) NOT NULL,
  etalase varchar(255) ,
  weight float NOT NULL,
  price int NOT NULL,
  rating float default 5,
  created_at timestamp NOT NULL default NOW(),
  updated_at timestamp default NOW()  
);

CREATE TABLE IF NOT EXISTS product_images (
  id serial PRIMARY KEY,
  product_id bigint NOT NULL,
  image_url varchar(255) NOT NULL,
  short_description varchar(50),
  created_at timestamp NOT NULL default NOW(),
  updated_at timestamp default NOW()  
);

CREATE TABLE IF NOT EXISTS product_reviews (
  id serial PRIMARY KEY,
  product_id bigint NOT NULL,
  rating int NOT NULL,
  comment varchar(255),
  created_at timestamp NOT NULL default NOW(),
  updated_at timestamp default NOW()  
);