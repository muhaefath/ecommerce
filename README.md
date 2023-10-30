# ecommerce

#### HOW to run the apps

##### 1. run $docker-compose up

##### 2. migrate the database with command below

###### $ migrate -source file://migration -database 'postgres://postgres:ecommerce@localhost:5432?sslmode=disable' up 2

##### 3. run the application with command below

###### go run main.go

##### 4. import postman collection and call the API

