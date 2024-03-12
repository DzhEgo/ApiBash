# Инструкция по запуску
## Загрузить архив с кодом и запустить на локальной машине или поднять его в докере. Открываем Postman или используем командну curl в терминале. Далее вводим localhost:{port} и отправляем данные. 

## В приложении присутствует 3 хэндлера с post, get и getById запросами. 
### Внешние пакеты, которые использовал: Gorm, gorilla/mux
### CommPostHandler - осуществляет добавление и применение bash скрипта.
### CommGetHandler - осуществляет вывод всех bash скриптов, которые находятся в базе, вывод и результаты выполнения скриптов.
### CommGetIdHandler - осуществляет вывод по id и то же самое, что и CommGetHandler.
### Формат json'а (request):
```json
  "Script":"echo Hello world!"
```
### Формат json'а (response):
```json
  "Script":"echo Hello world!"
  "Result":"Hello world!"
  "Executed":"true"
```
### Также, на выходе дополнительные параметры взятые из модели gorm.Model такие как: id, createdAt, updatedAt, deletedAt 
