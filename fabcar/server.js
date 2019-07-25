var query = require("./query.js");
var invoke = require("./create.js");
var express = require("express");
var fs = require("fs");
var cors = require('cors')
var bodyParser = require("body-parser")
var app = express();


app.use(cors())
app.use(bodyParser.urlencoded({
  extended: true
}));
app.use(express.json());

app.get('/queryCar', function(request, response) {
    response.writeHead(200, {"Content-Type":"text/html"});
    fs.readFile("html/queryCar.html", "utf-8", function(e, data){
        response.write(data);
        response.end();
    });
});

app.get('/createCar', function(request, response) {
    response.writeHead(200, {"Content-Type":"text/html"});
    fs.readFile("html/createCar.html", "utf-8", function(e, data){
        response.write(data);
        response.end();
    });
});

app.get('/changeOwner', function(request, response) {
    response.writeHead(200, {"Content-Type":"text/html"});
    fs.readFile("html/changeOwner.html", "utf-8", function(e, data){
        response.write(data);
        response.end();
    });
});

app.post('/queryMER', function(request, response) {
    var patientId = request.body.patientId;
    console.log(request.body)
    // response.write(request.body)
    console.log(patientId)
    query.queryMER(patientId).then((result) => {
             response.writeHead(200, {'Content-type': 'application/json'});
             console.log(result)
             if (result.length == 0){
                 result = "car not found!" 
             }
            //  response.write(JSON.stringify({patientId,result}));
            response.write(JSON.stringify(result[0]))
             response.end();
          });
});
app.post('/create', function(request, response) {
    args = request.body.args;
    query.create(args).then((result) => {
             response.writeHead(200, {'Content-type': 'application/json'});
             if (result.length == 0){
                 result = "car not found!" 
             }
             response.write(result);
             response.end();
          });
});
app.post('/invoke', function(request, response) {
    func = request.body.func
    console.log(func)

    if (func == 'createCar'){
        carID = request.body.carID;
        make = request.body.make;
        module = request.body.module;
        colour = request.body.colour;
        owner = request.body.owner;
        invoke.invokecc(func, [carID, make, module, colour, owner])
            .then((result) => {
                 response.writeHead(200, {'Content-type': 'application/json'});
                 response.write(result);
                 response.end();
            });
     } else if(func == 'changeCarOwner'){
         carID = request.body.carID;
         owner = request.body.newOwner;
         invoke.invokecc(func, [carID, owner])
             .then((result) => {
                  response.writeHead(200, {'Content-type': 'application/json'});
                  response.write(result);
                  response.end();
             });
     }
});
//设置跨域访问
app.all('*', function(req, res, next) {
    res.header("Access-Control-Allow-Origin", "*");
    res.header('Access-Control-Allow-Headers', 'Content-Type, Content-Length, Authorization, Accept, X-Requested-With , yourHeaderFeild');
    res.header("Access-Control-Allow-Methods","PUT,POST,GET,DELETE,OPTIONS");
    res.header("X-Powered-By",' 3.2.1')
    res.header("Content-Type", "application/json;charset=utf-8");
    next();
});
console.log("Listening on port 8080")
app.listen(8080)
