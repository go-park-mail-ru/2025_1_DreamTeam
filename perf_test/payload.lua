wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"
wrk.headers["User-Agent"] = "PostmanRuntime/7.43.4"
wrk.headers["Accept"] = "*/*"
wrk.headers["Postman-Token"] = "7cf9a615-3231-477b-aa2b-b308cc445a5e"
wrk.headers["Accept-Encoding"] = "gzip, deflate, br"
wrk.headers["Connection"] = "keep-alive"
wrk.headers["Cookie"] = "session_id=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjozOH0.gRvcEDXiCk8SmAZs0ixyw_UylT2xMdSlm4LDg2VbWAw"

wrk.body = [[
{
  "price": 0,
  "time_to_pass": 30,
  "title": "Testing DB",
  "description": "<h2>На курсе освоите</h2>\n<p>Ключевые концепции, стили и произведения двух главных литературных направлений XIX века</p>",
  "parts": [
    {
      "part_title": "Мир романтиков",
      "buckets": [
        {
          "bucket_title": "Бунтарские духи",
          "lessons": [
            {
              "lesson_type": "text",
              "lesson_title": "lesson",
              "lesson_value": "Байрон и Гюго: культ героя-одиночки. Ключевые произведения и цитаты"
            },
            {
              "lesson_type": "text",
              "lesson_title": "lesson",
              "lesson_value": "Концепция 'мировой скорби' в романтизме: от Вертера до Лермонтова"
            }
          ]
        }
      ]
    }
  ]
}
]]
