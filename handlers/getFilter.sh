curl -X POST http://localhost:5555/hotel \
     -H "Content-Type: application/json" \
     -d '[{"filter": "Mayerthorpe,Alberta,53.95015;-115.13547"},{"filter": "Taber,Alberta,49.78703;-112.14603"}]'