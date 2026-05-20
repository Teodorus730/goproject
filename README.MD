# Тесты


## Винда

### Рега
```
$body = @{                                                                      
name     = "PETYA"            
email    = "ivan@example.com" 
phone    = "79991234567"    
password = "secret_password"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/register" -Method Post -ContentType "application/json" -Body $body
```

### Логин
```
$loginBody = @{
    email    = "ivan@example.com"
    password = "secret_password"
} | ConvertTo-Json

$loginResponse = Invoke-RestMethod -Uri "http://localhost:8080/login" -Method Post -ContentType "application/json" -Body $loginBody 
```

### Создать заказ
```
$orderBody = @{
    tarif_id = 1
    amount   = 10
    address  = "SixSeven St. 67"
} | ConvertTo-Json

$headers = @{
    Authorization = "Bearer $($loginResponse.token)"
}

Invoke-RestMethod -Uri "http://localhost:8080/add_order" -Method Post -ContentType "application/json" -Body $orderBody -Headers $headers   
```

### Все заказы пользователя
```
$headers = @{
    Authorization = "Bearer $($loginResponse.token)"
}

Invoke-RestMethod -Uri "http://localhost:8080/get_orders" -Method Get -Headers $headers   
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

# Создание тарифов

```
docker exec -it postgres psql -U postgres -d database -c "INSERT INTO tariffes (price, liters, name) VALUES (100, 10, 'Desyat litrov')"
docker exec -it postgres psql -U postgres -d database -c "INSERT INTO tariffes (price, liters, name) VALUES (200, 20, 'Dvatsat litrov')"
```