var query = require("./queryExport.js");
var invoke = require("./invokeExport.js");
var express = require("express");
var fs = require("fs");
var bodyParser = require("body-parser")

var app = express();
app.use(bodyParser.urlencoded({
  extended: true
}));

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

app.post('/queryCar', function(request, response) {
    car = request.body.car;
    query.queryCAR(car)
         .then((result) => {
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

console.log("Listening on port 8080")
app.listen(8080)
