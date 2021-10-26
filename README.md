# REST
Examples of REST web server implementation in Go with different storages.


Already DONE:

web server
1. Standard library net/http web server - pkg stdlib-http;
2. Library gorilla/mux web server - pkg gorilla;
3. Library gin web server - pkg gin-gonic;
4. Library fasthttp web server - pkg fasth;

storage
1. in-memory storage using map - pkg inmemory;

data model - pkg models. 

config - pkg config, file conig.yaml. 

test queries - file testURL.txt.



TODO:

Middleware

OpenAPI && Swagger

Graceful shutdown

web server
1. Library fiber web server;
2. Library beego web server;
3. Library echo web server;
4. Library kit web server;

storage
1. PostgreSQL;
2. MySQL;
3. MongoDB;

authentication
1. HTTPS/TLS;
2. JWT;
3. Cookies;
4. OAuth 2.0;

Unit tests

Benchmarks
