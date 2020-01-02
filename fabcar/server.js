var query = require("./queryCase.js");
var querycaseinfo = require("./queryCaseInfo.js");
var querypatientinfo = require("./queryPatientInfo.js");
var Create = require("./create.js");
var invoke = require("./invoke.js");
var addreserveinfo = require('./addreserveinfo.js');
var queryreserve = require('./queryyuyue.js');
var queryreserveinfo = require('./queryreserveinfo.js');
var querydoc = require('./querydoc.js');
var querydoctorinfo = require('./querydoctorinfo');
var express = require("express");
var fs = require("fs");
var cors = require('cors');
var bodyParser = require("body-parser");
var app = express();
var querymedicine = require("./querymedicine.js")
//var sendmessage = require('./message')
const Core = require('@alicloud/pop-core');


app.use(cors())
app.use(bodyParser.urlencoded({
    extended: true
}));
app.use(express.json());

app.get('/queryCar', function (request, response) {
    response.writeHead(200, { "Content-Type": "text/html" });
    fs.readFile("html/queryCar.html", "utf-8", function (e, data) {
        response.write(data);
        response.end();
    });
});

app.get('/createCar', function (request, response) {
    response.writeHead(200, { "Content-Type": "text/html" });
    fs.readFile("html/createCar.html", "utf-8", function (e, data) {
        response.write(data);
        response.end();
    });
});

app.get('/changeOwner', function (request, response) {
    response.writeHead(200, { "Content-Type": "text/html" });
    fs.readFile("html/changeOwner.html", "utf-8", function (e, data) {
        response.write(data);
        response.end();
    });
});

app.post('/queryCase', function (request, response) {
    var patientId = request.body.patientId;
    console.log(request.body)
    // response.write(request.body)
    console.log(patientId)
    query.queryCase(patientId).then((result) => {
        response.writeHead(200, { 'Content-type': 'application/json' });
        //  console.log(result)
        if (result.length == 0) {
            result = "car not found!"
        }
        // response.write(JSON.stringify({patientId,result}));
        var str = "" + result
        // console.log(str)
        // console.log(JSON.stringify(result))
        response.write(str)
        response.end();
    });
});
app.post('/querymedicine', function (request, response) {
    var patientId = request.body.patientId;
    var medicineid = request.body.medicineid;

    querymedicine.querymedicine(patientId, medicineid).then((result) => {
        response.writeHead(200, { 'Content-type': 'application/json' });
        //  console.log(result)
        if (result.length == 0) {
            result = "car not found!"
        }
        // response.write(JSON.stringify({patientId,result}));
        var str = "" + result
        // console.log(str)
        // console.log(JSON.stringify(result))
        response.write(str)
        response.end();
    });
});
app.post('/querydoc', function (request, response) {
    var hospitalname = request.body.hospitalname;
    var department = request.body.department;
    querydoc.querydoc(hospitalname, department).then((result) => {
        response.writeHead(200, { 'Content-type': 'application/json' });
        //  console.log(result)
        if (result.length == 0) {
            result = "car not found!"
        }
        // response.write(JSON.stringify({patientId,result}));
        var str = "" + result
        // console.log(str)
        // console.log(JSON.stringify(result))
        response.write(str)
        response.end();
    });
});
app.post('/querydoctorinfo', function (request, response) {
    var doctorId = request.body.doctorId;
    // console.log(request.body)
    // response.write(request.body)
    // console.log(patientId)
    querydoctorinfo.querydoctorinfo(doctorId).then((result) => {
        response.writeHead(200, { 'Content-type': 'application/json' });
        //  console.log(result)
        if (result.length == 0) {
            result = "car not found!"
        }
        //  response.write(JSON.stringify({patientId,result}));
        var str = "" + result
        response.write(str)
        response.end();
    });
});
app.post('/queryCaseInfo', function (request, response) {
    var patientId = request.body.patientId;
    var complaintid = request.body.complaintid;
    // console.log(request.body)
    // response.write(request.body)
    // console.log(patientId)
    querycaseinfo.queryCaseInfo(patientId, complaintid).then((result) => {
        response.writeHead(200, { 'Content-type': 'application/json' });
        //  console.log(result)
        if (result.length == 0) {
            result = "car not found!"
        }
        //  response.write(JSON.stringify({patientId,result}));
        var str = "" + result
        response.write(str)
        response.end();
    });
});
app.post('/queryPatientInfo', function (request, response) {
    var patientId = request.body.patientId;
    console.log(request.body)
    // response.write(request.body)
    console.log(patientId)
    querypatientinfo.queryPatientInfo(patientId).then((result) => {
        response.writeHead(200, { 'Content-type': 'application/json' });
        console.log(result)
        if (result.length == 0) {
            result = "car not found!"
        }
        //  response.write(JSON.stringify({patientId,result}));
        var str = "" + result
        response.write(str)
        response.end();
    });
});
app.post('/queryDoctorInfo', function (request, response) {
    var doctorId = request.body.doctorId;
    console.log(request.body)
    // response.write(request.body)
    //console.log(patientId)
    querypatientinfo.queryPatientInfo(doctorId).then((result) => {
        response.writeHead(200, { 'Content-type': 'application/json' });
        console.log(result)
        if (result.length == 0) {
            result = "car not found!"
        }
        //  response.write(JSON.stringify({patientId,result}));
        var str = "" + result
        response.write(str)
        response.end();
    });
});
app.post('/addreserveinfo', function (request, response) {
    patientId = request.body.patientId;
    reserveid = request.body.reserveid;
    hospitalname = request.body.hospitalname;
    department = request.body.department;
    doctorid = request.body.doctorid;
    reserverdate = request.body.reserverdate;
    reservertime = request.body.reservertime;
    reserverstate = request.body.reserverstate;
    addreserveinfo.addreserveinfo([patientId, reserveid, hospitalname, department, doctorid, reserverdate, reservertime, reserverstate]).then((result) => {
        response.writeHead(200, { 'Content-type': 'application/json' });
        if (result.length == 0) {
            result = "car not found!"
        }
        console.log(result)
        response.write(result);
        response.end();
    })
})
app.post('/sendmessage', function (request, response) {
    var number = request.body.number;
    var code = request.body.code
    console.log(request.body)
    // response.write(request.body)
    //  console.log(patientId)

    var client = new Core({
        accessKeyId: 'LTAI4FdFxdWgYzyW7689Fahi',
        accessKeySecret: 'cfz9yLqo5ArrNli1tvyty0V1cduaAJ',
        endpoint: 'https://dysmsapi.aliyuncs.com',
        apiVersion: '2017-05-25'
    });

    var params = {
        "RegionId": "cn-hangzhou",
        "PhoneNumbers": number,
        "SignName": "Ecase",
        "TemplateCode": "SMS_172887364",
        "TemplateParam": code
    }
    console.log(params)
    var requestOption = {
        method: 'POST'
    };

    client.request('SendSms', params, requestOption).then((result) => {
        console.log(result)
     var str =JSON.stringify(result);
        response.write(str);
        response.end();
    }, function (err) {
        console.log(err)
    })
    
});
app.post('/create', function (request, response) {
    patientId = request.body.patientid;
    complaintid = request.body.complainid;
    createdoctorid = request.body.createdoctorid;
    department = request.body.department;
    time = request.body.time;
    medicalType = request.body.medicalType;
    symptoms = request.body.symptoms;
    conclusion = request.body.conclusion;
    presenter = request.body.presenter;
    diseasesOnceSuffered = request.body.diseasesOnceSuffered;

    Create.Create([createdoctorid, patientId, complaintid, department, time, medicalType, symptoms, conclusion, presenter, diseasesOnceSuffered]).then((result) => {
        response.writeHead(200, { 'Content-type': 'application/json' });
        if (result.length == 0) {
            result = "car not found!"
        }
        var str = "" + result;
        response.write(str);
        response.end();
    });
});
app.post('/queryreserve', function (request, response) {
    patientId = request.body.patientId;
    queryreserve.queryreserve(patientId).then((result) => {
        response.writeHead(200, { 'Content-type': 'application/json' });
        var str = "" + result;
        response.write(str);
        response.end();
    });
});
app.post('/queryreserveinfo', function (request, response) {
    patientId = request.body.patientId;
    reserveid = request.body.reserveid;
    queryreserveinfo.queryreserveinfo(patientId, reserveid).then((result) => {
        response.writeHead(200, { 'Content-type': 'application/json' });
        if (result.length == 0) {
            result = "car not found!"
        }
        var str = "" + result;
        response.write(str);
        response.end();
    });
});
app.post('/invoke', function (request, response) {
    func = request.body.func;
    console.log(func);

    if (func == 'addReserveInfo') {
        patientId = request.body.patientId;
        reserveid = request.body.reserveid;
        hospitalname = request.body.hospitalname;
        department = request.body.department;
        doctorid = request.body.doctorid;
        reserverdate = request.body.reserverdate;
        reservertime = request.body.reservertime;
        reserverstate = request.body.reserverstate;
        invoke.invokecc(func, [patientId, reserveid, hospitalname, department, doctorid, reserverdate, reservertime, reserverstate])
            .then((result) => {
                console.log(result);
                response.writeHead(200, { 'Content-type': 'application/json' });
                var str = "" + result;
                response.write(str);
                response.end();
            });
    } else if (func == 'changeReserverState') {
        patientId = request.body.patientId;
        reserveid = request.body.reserveid;
        invoke.invokecc(func, [patientId, reserveid])
            .then((result) => {
                response.writeHead(200, { 'Content-type': 'application/json' });
                var str = "" + result;
                response.write(str);
                response.end();
            });
    } else if (func == 'patientregister') {

        username = request.body.username;
        usertype = request.body.usertype;
        password = request.body.password;
        patientid = request.body.patientid;
        mobliephone = request.body.mobliephone;
        rsapublic = request.body.rsapublic;
        invoke.invokecc(func, [username, usertype, password, patientid, mobliephone, rsapublic])
            .then((result) => {
                response.writeHead(200, { 'Content-type': 'application/json' });
                var str = "" + result;
                // var str = JSON.stringify(result)
                response.write(str);
                response.end();
            });
    } else if (func == 'patientLogin') {
        username = request.body.username;
        password = request.body.password;
        mobliephone = request.body.mobliephone;
        invoke.invokecc(func, [username, password, mobliephone])
            .then((result) => {
                response.writeHead(200, { 'Content-type': 'application/json' });
                var str = "" + result;
                //  str = JSON.stringify(result)
                response.write(str);
                response.end();
            });
    } else if (func == 'RSAPublicEncryptoAES') {
        username = request.body.username;
        invoke.invokecc(func, [username])
            .then((result) => {
                response.writeHead(200, { 'Content-type': 'application/json' });
                // var str = "" + result;
                //   str = JSON.stringify(result)
                response.write(result);
                response.end();
            });
    } else if (func == 'createPatient') {
        username = request.body.username;
        name = request.body.name;
        patientId = request.body.patientId;
        time = request.body.time;
        gender = request.body.gender;
        birthday = request.body.birthday;
        nation = request.body.nation;
        homeAddress = request.body.homeAddress;
        marriagecondition = request.body.marriagecondition;
        invoke.invokecc(func, [username, patientId, time, name, gender, birthday, nation, homeAddress, marriagecondition])
            .then((result) => {
                response.writeHead(200, { 'Content-type': 'application/json' });
                var str = "" + result;
                // var str = JSON.stringify(result)
                response.write(str);
                response.end();
            });

    } else if (func == 'aestest') {
        flag = request.body.flag;
        username = request.body.username;
        str = request.body.str;
        invoke.invokecc(func, [flag, username, str])
            .then((result) => {
                response.writeHead(200, { 'Content-type': 'application/json' });
                var str = "" + result;
                // var str = JSON.stringify(result)
                response.write(str);
                response.end();
            });
    } else if (func == 'changePaientBaseinfo') {
        patientId = request.body.patientId;
        homeAddress = request.body.homeAddress;
        marriageCondition = request.body.marriageCondition;
        invoke.invokecc(func, [patientId, homeAddress, marriageCondition])
            .then((result) => {
                response.writeHead(200, { 'Content-type': 'application/json' });
                var str = "" + result;
                // var str = JSON.stringify(result)
                response.write(str);
                response.end();
            });
    } else if (func == 'checkPermission') {
        patientId = request.body.patientId;
        doctorId = request.body.doctorId;
        invoke.invokecc(func, [doctorId, patientId])
            .then((result) => {
                response.writeHead(200, { 'Content-type': 'application/json' });
                var str = "" + result;
                // var str = JSON.stringify(result)
                response.write(str);
                response.end();
            });
    } else if (func == 'doctorPrescribe') {
        patientId = request.body.patientId;
        complainid = request.body.complainid;
        medicineid = request.body.medicineid;
        medicinename = request.body.medicinename;
        specification = request.body.specification;
        directions = request.body.directions;
        remark = request.body.remark;
        invoke.invokecc(func, [patientId, complainid, medicineid, medicinename, specification, directions, remark])
            .then((result) => {
                response.writeHead(200, { 'Content-type': 'application/json' });
                var str = "" + result;
                // var str = JSON.stringify(result)
                response.write(str);
                response.end();
            });
    } else if (func == 'transferPermission') {
        doctorId = request.body.doctorId;
        patientId = request.body.patientId;
        invoke.invokecc(func, [doctorId, patientId])
            .then((result) => {
                response.writeHead(200, { 'Content-type': 'application/json' });
                var str = "" + result;
                // var str = JSON.stringify(result)
                response.write(str);
                response.end();
            });
    } else if (func == 'doctorLogin') {
        username = request.body.username;
        password = request.body.password;
        invoke.invokecc(func, [username, password])
            .then((result) => {
                response.writeHead(200, { 'Content-type': 'application/json' });
                var str = "" + result;
                // var str = JSON.stringify(result)
                response.write(str);
                response.end();
            });
    }
    else if (func == 'createDoctor') {
        doctorId = request.body.doctorId;
        name = request.body.name;
        hospitalname = request.body.hospitalname;
        role = request.body.role;
        department = request.body.department;
        price = request.body.price;
        doctordate = request.body.doctordate;
        doctortime = request.body.doctortime;
        invoke.invokecc(func, [doctorId, name, hospitalname, role, department, price, doctordate, doctortime])
            .then((result) => {
                response.writeHead(200, { 'Content-type': 'application/json' });
                var str = "" + result;
                // var str = JSON.stringify(result)
                response.write(str);
                response.end();
            });
    }
});
//设置跨域访问
app.all('*', function (req, res, next) {
    res.header("Access-Control-Allow-Origin", "*");
    res.header('Access-Control-Allow-Headers', 'Content-Type, Content-Length, Authorization, Accept, X-Requested-With , yourHeaderFeild');
    res.header("Access-Control-Allow-Methods", "PUT,POST,GET,DELETE,OPTIONS");
    res.header("X-Powered-By", ' 3.2.1')
    res.header("Content-Type", "application/json;charset=utf-8");
    next();
});
console.log("Listening on port 8080")
app.listen(8080);
