Это простое консольное приложение на Go для управления пулом воркеров. Воркеры могут добавляться, удаляться и выполнять задачи в фоновом режиме.

## Начало работы

1.  **Клонирование репозитория:**
    ```sh
    git clone https://github.com/nongrata2/workerPool
    cd workerPool
    ```
    
2.  **Запуск приложения:**
  ```sh
  go run cmd/main.go
  ```

## Команды приложения

| Команда           | Описание                                                          | Пример использования      |
| :---------------- | :---------------------------------------------------------------- | :------------------------ |
| `add`             | Добавляет нового воркера в пул с автоматически присвоенным ID.    | `add`                     |
| `addid [id]`      | Добавляет нового воркера с указанным ID.                          | `addid 5`                 |
| `ls`              | Отображает список активных воркеров.                              | `ls`                      |
| `remove`          | Удаляет воркера с наибольшим ID из пула.                         | `remove`                  |
| `removeid [id]`   | Удаляет воркера с указанным ID из пула.                          | `removeid 3`              |
| `process [task]`  | Отправляет задачу в очередь воркеров. Воркеры будут обрабатывать ее. | `process mytask`   |
| `help`            | Выводит список всех доступных команд и их описание.               | `help`                    |
| `exit`            | Завершает работу приложения и останавливает всех воркеров.       | `exit`                    |

Для запуска тестов необходимо использовать
```sh
go test ./... -v 
```