# REST
Examples of REST web server implementation in Go with different storages.


### Already DONE:

Web server
1. Standard library net/http web server - pkg stdlib-http;
2. Library gorilla/mux web server - pkg gorilla;
3. Library gin web server - pkg gin-gonic;
4. Library fasthttp web server - pkg fasth;

Storage
1. in-memory storage using map - pkg inmemory;

Other features:
- Data model - pkg models. 
- Config - pkg config, file conig.yaml. 
- Test queries - file testURL.txt.
- Middleware - pkg middleware. 

### TODO:

Web server
1. Library fiber web server;
2. Library beego web server;
3. Library echo web server;
4. Library kit web server;

Storage
1. PostgreSQL;
2. MySQL;
3. MongoDB;

Authentication
1. HTTPS/TLS;
2. JWT;
3. Cookies;
4. OAuth 2.0;

Other features:
- OpenAPI && Swagger
- Graceful shutdown
- Unit tests
- Benchmarks
