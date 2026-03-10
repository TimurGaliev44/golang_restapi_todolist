How to run?
1. Install
git clone git@github.com:TimurGaliev44/golang_restapi_todolist.git
2. Download dependencies
go mod download
3. Run docker 
docker-compose up

Now you can test this project, just send HTTP request on 5070
Example: curl -i -X POST http://localhost:5070/register   -H "Content-Type: application/json"   -d '{"username":"test","password":"test123"}'
