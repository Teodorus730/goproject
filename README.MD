# Тесты


## Винда

```
$body = @{                                                                      
name     = "PETYA"            
email    = "ivan@example.com" 
phone    = "79991234567"    
password = "secret_password"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/register" -Method Post -ContentType "application/json" -Body $body
```


## Нормальная ОС

```
curl -X POST http://localhost:8080/register \
     -H "Content-Type: application/json" \
     -d '{
       "name": "Ivan",
       "email": "ivan@example.com",
       "phone": "79991234567",
       "password": "secret_password"
     }'
```


# Тест в БД

```
 docker exec -it postgres psql -U postgres -d database -c "SELECT * FROM users;"
```