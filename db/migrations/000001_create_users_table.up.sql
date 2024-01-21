CREATE TABLE users(
   id serial PRIMARY KEY,
   name VARCHAR (50) NOT NULL,
   surname VARCHAR (50) NOT NULL,
   patronymic VARCHAR (50) NOT NULL,
	age INTEGER NOT NULL,
	gender VARCHAR (50) NOT NULL,
	nationality VARCHAR (50) NOT NULL
);