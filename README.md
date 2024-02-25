# bmstu-ppo

# Тема: Музыкальный сервис

### Идея проекта
Создать приложение, которое предоставляет пользователю возможность авторизоваться и прослушивать треки.
Понравившиеся треки можно добавить в избранное.  
Приложение позволяет просматривать события конкретного исполнителя.

### Описание предметной области
Предметная область --- приложение музыкального сервиса, включающее: регистрацию пользователей, добавление новых музыкальных композиций, формирование списка избранных треков.


### Анализ аналогичных решений

| Решение  | Избранное | События исполнителя | Похожие треки | Доступно в России |
|----------|----------|----------| -- | - |
| Яндекс Музыка    | +   | - | + | + |
| VK музыка    | +   | -   | + | + |
| Spotify    | +   | + | + | - |
| Разрабатываемое приложение | + | + | + | + |


### Целесообразность и актуальность проекта
Музыка занимает важное место в жизни людей, поэтому музыкальные сервисы, позволяющие прослушивать треки и выбирать из них понравившиеся востребованы всегда.

### Акторы
- Обычный пользователь (слушатель) --- роль, которая может только прослушивать треки и не может добавлять новые
- Музыкант --- роль, которая позволяет только добавлять новые треки от своего имени, те не может влиять на треки других ролей
- Администратор --- роль, которая позволяет добавлять и удалять новых пользователей, музыкантов и редактировать все треки.

### Use-Case диаграмма
![usecase](./docs/usecase.png)

### 