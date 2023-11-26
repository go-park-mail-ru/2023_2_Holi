hey -z 5m -q 200 -m POST http://localhost:3001/api/v1/auth/register &&
hey -z 5m -q 400 -m POST http://localhost:3001/api/v1/auth/login &&
hey -z 5m -q 200 -m POST http://localhost:3001/api/v1/auth/logout &&
hey -z 5m -q 200 -m POST http://localhost:3001/api/v1/auth/check
