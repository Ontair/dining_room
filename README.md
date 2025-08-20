# dining_room



curl -X POST http://localhost:8080/dish \
  -H "Content-Type: application/json" \
  -d '{
    "name": "пюре",
    "price": "100",
    "description": "картоха"
  }'