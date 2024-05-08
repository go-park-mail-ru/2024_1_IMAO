# 2024_1_IMAO
Backend репозиторий команды IMAO

### Разработчики:

Алексей Оглоблин: https://github.com/AlexeiLDD

Илья Чугин: https://github.com/IlyaChgn

Марат Камалов: https://github.com/MaratKamalovPD

Олег Жданов: https://github.com/telephonist95

### Репозиторий фронтенда:
https://github.com/frontend-park-mail-ru/2024_1_IMAO

### Ссылка на Фигму:
https://www.figma.com/file/QP3qZTavTZYL8aOlzWzkkl/IMAO-(%D0%AE%D0%BB%D0%B0)

### Деплой приложения:
http://www.vol-4-ok.ru:8008

### Запуск локально:
`go run cmd/app/main.go`

### Тестирование
```
go test --cover  ./...
```
или
```
mkdir -p bin && go test -coverprofile=bin/cover.out ./internal/... && go tool cover -html=bin/cover.out -o=bin/cover.html && go tool cover --func bin/cover.out
```


## Документация
Можно посмотреть всю информацию в docs/swagger.yaml

### Сгенерировать swagger документацию
`swag init -g cmd/app/main.go`
