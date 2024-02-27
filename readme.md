# 怎么运行mongodb

```bash
docker run -d --name mongodb -v ./mongo:/data/db -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=password -p 27017:27017 --rm mongo


JWT_SECRET=eUbP9shywUygMx7u MONGO_URI="mongodb://admin:password@localhost:27017/test?authSource=admin" MONGO_DATABASE=test go run main.go

# mongoimport data to MongoDB
mongoimport --username admin --password password --authenticationDatabase admin --db test --collection recipes --file recipes.json --jsonArray
```

## 注意事项

1. 使用`token.SignedString`时，参数必须是`[]byte`