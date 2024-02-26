# 怎么运行mongodb

```bash
docker run -d --name mongodb -v ./mongo:/data/db -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=password -p 27017:27017 --rm mongo


MONGO_URI="mongodb://admin:password@localhost:27017/test?authSource=admin" MONGO_DATABASE=test go run main.go

# mongoimport data to MongoDB
mongoimport --username admin --password password --authenticationDatabase admin --db test --collection recipes --file recipes.json --jsonArray
```