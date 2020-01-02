package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	mathrand "math/rand"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type merChaincode struct {
}

//帐号信息
type User struct {
	ObjectType   string `json:"docType"`      //field for couchdb
	UserNmae     string `json:"userNmae"`     //用户名
	UserType     string `json:"userType"`     //用户的类型
	Password     string `json:"password"`     //用户的password
	IDcardNumber string `json:"idcardNumber"` //用户的身份证号
	MobilePhone  string `json:"mobilePhone"`  //用户的手机号
	AESKey       string `json:"aesKey"`       //用户的aes密钥
	RSAPublic    string `json:"rsaPublic"`    //用户的rsa公钥
}

//患者的所有信息
type Patient struct {
	ObjectType        string           `json:"docType"`      //field for couchdb
	IdCardNumber      string           `json:"idcardnumber"` //患者身份证号
	BaseInfo          PatientBaseInfo  `json:"baseInfo"`
	MedicalContents   [100]Complaint   `json:"medicalContents"`   //病历数组
	ReserveInfo       [100]Reservation `json:"reserveInfo"`       //预约信息数组
	DoctorIDcollect   [100]string      `json:"doctorIDcollect"`   //医生的ID数组
	MedicalContentNum int              `json:"medicalContentNum"` //病历的总个数
	ReserveInfoNum    int              `json:"reserveInfoNum"`    //预约信息的总个数
	DoctorIDNum       int              `json:"doctorIDNum"`       //对应医生的总个数
}

// 患者基本信息
type PatientBaseInfo struct {
	IdCardNumber      string `json:"idcardnumber"` //患者身份证号
	CreateTime        string `json:"createtime"`   //信息创建时间
	PatientName       string `json:"patientName"`
	Gender            string `json:"gender"` //性别m为男
	Birthday          string `json:"birthday"`
	Nation            string `json:"nation"`            //民族
	HomeAddress       string `json:"homeAddress"`       //家庭住址
	MarriageCondition string `json:"marriageCondition"` //婚姻情况
}

//患者病历信息
type Complaint struct {
	ComplaintId          string      `json:"complaintId"`          //患者该条病历的 ID
	CreateDoctorId       string      `json:"createDoctorId"`       //该条病历创建的医生id
	Department           string      `json:"department"`           //科室
	MedicalCreateTime    string      `json:"medicalCreate"`        //该条病历的创建时间
	MedicalType          string      `json:"medicalType"`          //患者疾病类型
	MainSymptoms         string      `json:"mainSymptoms"`         //患者主诉
	Conclusion           string      `json:"conclusion"`           //医生的诊断结果
	Presenter            string      `json:"presenter"`            //陈述者
	DiseasesOnceSuffered string      `json:"diseasesOnceSuffered"` //患者曾患病
	MedicineContentNum   int         `json:"medicineNum"`          //该条病历的所有药品数量
	MedicineIDContent    [100]string `json:"medicineIDContent"`    //病历所对应的全部药品
	//PatientPersonalDesctiption  string   `json:"patientPersonalDesctiption"`    //患者个人史***************
	//DetailSymptoms    string    `json:"detailSymptoms"`      //患者症状详细信息   *************************
}

//药品信息
type Medicine struct {
	MedicineId    string `json:"medicineId"`
	MedicineName  string `json:"medicineName"`  //药品名称
	Specification string `json:"specification"` //药品规格
	Directions    string `json:"directions"`    //药品用法
	Remark        string `json:"remark"`        //药品备注
}

//预约信息
type Reservation struct {
	ReserveID     string `json:"reserveID"`     //预约ID
	HospitalName  string `json:"hospitalName"`  //医院名称
	Department    string `json:"department"`    //预约科室
	DoctorId      string `json:"doctorId"`      //医生ID
	ReserverDate  string `json:"reserverDate"`  //挂号日期
	ReserverTime  string `json:"reserverTime"`  //挂号时间 上/下午
	ReserverState string `json:"reserverState"` //预约状态
}

//医生信息
type Doctor struct {
	ObjectType       string      `json:"docType"`  //field for couchdb
	DoctorId         string      `json:"doctorId"` //医生ID
	DoctorName       string      `json:"doctorName"`
	HospitalName     string      `json:"hospitalname"`     //医生所属医院的名字
	Role             string      `json:"role"`             //科室角色 主任医师等
	Department       string      `json:"department"`       //所在科室
	Price            float64     `json:"price"`            //挂号价格
	DoctorDate       string      `json:"doctorDate"`       //坐诊日期
	DoctorTime       string      `json:"doctorTime"`       //坐诊时间 上/下午
	PatientIDcollect [100]string `json:"patientIDcollect"` //患者的ID数组
	PatientIDNum     int         `json:"patientIDNum"`     //与该医生关联的病人id的个数
}

//患者症状详细描述信息
type DetailSymptomsContent struct {
	OnesetTimeAndPossiblePathogeny   string `json:"onesetTimeAndPossiblePathogeny"`   //发病时间和可能原因
	MainSymptomsElaborateDescription string `json:"mainSymptomsElaborateDescription"` //主要症状的详细描述
	SimultaneousPhenomenon           string `json:"simultaneousPhenomenon"`           //伴随症状
	OtherDiseases                    string `json:"otherDiseases"`                    //患者患有的其他疾病
	GeneralConditions                string `json:"generalConditions"`                //发病以来的一般情况
}

//患者个人史信息
type PatientPersonalDescription struct {
	RecentAreaAndDate           string `json:"recentAreaAndDate"`           //患者最近去过的地点以及时间
	HabitsAndCustoms            string `json:"habitsAndCustoms"`            //患者起居习惯等
	Occupation                  string `json:"occupation"`                  //患者的职业和工作环境
	VisitProstitutes            string `json:"visitProstitutes"`            //患者有无冶游史，以及发病时间
	GrowthAndDevelopmentHistory string `json:"growthAndDevelopmentHistory"` //患者生长发育史
	MaritalStaus                string `json:"maritalStaus"`                //患者婚姻情况
	Menstruation                string `json:"menstruation"`                //女性患者月经情况
}

//链代码初始化函数
func (t *merChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Patient Init Success")
	return shim.Success(nil)
}

//Invoke 函数
func (t *merChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	//根据function 参数值调用相应的处理函数
	if function == "patientregister" {
		return t.patientregister(stub, args)
	} else if function == "patientLogin" {
		return t.patientLogin(stub, args)
	} else if function == "createPatient" {
		return t.createPatient(stub, args)
	} else if function == "queryPaientBaseinfo" {
		return t.queryPaientBaseinfo(stub, args)
	} else if function == "changePaientBaseinfo" {
		return t.changePaientBaseinfo(stub, args)
	} else if function == "addMedicalContent" {
		return t.addMedicalContent(stub, args)
	} else if function == "queryMedicalByID" {
		return t.queryMedicalByID(stub, args)
	} else if function == "queryMedicalNum" {
		return t.queryMedicalNum(stub, args)
	} else if function == "addReserveInfo" {
		return t.addReserveInfo(stub, args)
	} else if function == "queryReserverInfoNum" {
		return t.queryReserverInfoNum(stub, args)
	} else if function == "queryReserverInfoByID" {
		return t.queryReserverInfoByID(stub, args)
	} else if function == "createDoctor" {
		return t.createDoctor(stub, args)
	} else if function == "createDoctorAuto" {
		return t.createDoctorAuto(stub, args)
	} else if function == "queryDepartmentByDoctorID" {
		return t.queryDepartmentByDoctorID(stub, args)
	} else if function == "deleteMedicalContent" { //测试用途
		return t.deleteMedicalContent(stub, args)
	} else if function == "queryDoctorByHospitalDepartment" {
		return t.queryDoctorByHospitalDepartment(stub, args)
	} else if function == "deleteReserverInfo" { //测试用途
		return t.deleteReserverInfo(stub, args)
	} else if function == "changeReserverState" {
		return t.changeReserverState(stub, args)
	} else if function == "queryDoctorInfoByID" {
		return t.queryDoctorInfoByID(stub, args)
	} else if function == "RSAPublicEncryptoAES" {
		return t.RSAPublicEncryptoAES(stub, args)
	} else if function == "checkPermission" {
		return t.checkPermission(stub, args)
	} else if function == "aestest" {
		return t.aestest(stub, args)
	} else if function == "doctorLogin" {
		return t.doctorLogin(stub, args)
	} else if function == "doctorPrescribe" {
		return t.doctorPrescribe(stub, args)
	} else if function == "transferPermission" {
		return t.transferPermission(stub, args)
	} else if function == "queryMedicineByNum" {
		return t.queryMedicineByNum(stub, args)
	} else if function == "getHistoryForPatient" {
		return t.getHistoryForPatient(stub, args)
	}
	return shim.Error("没找到对应方法~")
}

/***********************************************************************
     加密模块 各个加密解密函数
************************************************************************/
const iv = "1234567887654321"

//aes加密
func AesEncrypt(encodeStr string, k string) string {
	key := []byte(k)

	encodeBytes := []byte(encodeStr)
	//根据key 生成密文
	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()
	encodeBytes = PKCS5Padding(encodeBytes, blockSize)

	blockMode := cipher.NewCBCEncrypter(block, []byte(iv))
	crypted := make([]byte, len(encodeBytes))
	blockMode.CryptBlocks(crypted, encodeBytes)
	return base64.StdEncoding.EncodeToString(crypted)
}

//aes解密
func AesDecrypt(decodeStr string, k string) string {
	//先解密base64
	key := []byte(k)
	decodeBytes, err := base64.StdEncoding.DecodeString(decodeStr)
	if err != nil {
		return "error!"
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "error!!!"
	}
	blockMode := cipher.NewCBCDecrypter(block, []byte(iv))
	origData := make([]byte, len(decodeBytes))
	blockMode.CryptBlocks(origData, decodeBytes)
	origData = PKCS5UnPadding(origData)
	return string(origData)
}

//补码
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	//填充
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(ciphertext, padtext...)
}

//去码
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//rsa 加密
func RsaEncrypt(origData []byte, publicKey []byte) ([]byte, error) {
	//解密pem格式的公钥
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	//加密
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

func (t *merChaincode) aestest(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//flag := args[0]
	username := args[1]
	str := args[2]
	userAsBytes, err := stub.GetState(username)
	if err != nil {
		return shim.Error(err.Error())
	}
	userTemp := User{}
	json.Unmarshal(userAsBytes, &userTemp)
	aeskey := userTemp.AESKey
	//es := AesEncrypt(str, aeskey)
	ds := AesDecrypt(str, aeskey)
	// if flag == "0"{
	// 	return shim.Success([]byte(es))
	// }
	return shim.Success([]byte(ds))

}

/*******************************************************************************/

var RSApublicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDfw1/P15GQzGGYvNwVmXIGGxea
8Pb2wJcF7ZW7tmFdLSjOItn9kvUsbQgS5yxx+f2sAv1ocxbPTsFdRc6yUTJdeQol
DOkEzNP0B8XKm+Lxy4giwwR5LJQTANkqe4w/d9u129bRhTu/SUzSUIr65zZ/s6TU
GQD6QzKY1Y8xS+FoQQIDAQAB
-----END PUBLIC KEY-----
`)

//用rsa公钥加密用户的aes密钥
func (t *merChaincode) RSAPublicEncryptoAES(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	username := args[0]

	userAsBytes, err := stub.GetState(username)
	if err != nil {
		return shim.Error(err.Error())
	}

	userTemp := User{}
	json.Unmarshal(userAsBytes, &userTemp)

	aeskey := []byte(userTemp.AESKey)
	encryptoedAES, err := RsaEncrypt(aeskey, RSApublicKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	res := []byte(base64.StdEncoding.EncodeToString(encryptoedAES))

	return shim.Success(res)
}

//用医生rsa公钥加密患者的aes密钥
func (t *merChaincode) DoctorRSAPublicEncryptoAES(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	username := args[0]

	userAsBytes, err := stub.GetState(username)
	if err != nil {
		return shim.Error(err.Error())
	}

	userTemp := User{}
	json.Unmarshal(userAsBytes, &userTemp)

	aeskey := []byte(userTemp.AESKey)
	encryptoedAES, err := RsaEncrypt(aeskey, RSApublicKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	res := []byte(base64.StdEncoding.EncodeToString(encryptoedAES))

	return shim.Success(res)
}

//随机生成 16 位的字符串
func Krand(size int, kind int) []byte {
	ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	is_all := kind > 2 || kind < 0
	mathrand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if is_all { // random ikind
			ikind = mathrand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + mathrand.Intn(scope))
	}
	return result
}

//患者在患者系统注册
func (t *merChaincode) patientregister(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	username := args[0]
	//查询系统中有无该帐号
	userAsBytes, err := stub.GetState(username)
	if err != nil {
		return shim.Error(err.Error())
	}

	if userAsBytes != nil {
		s := "this user is already existed !!!!"
		ss := []byte(s)
		return shim.Success(ss)
	}
	objectType := "user"
	usertype := args[1]
	password := args[2]
	idcardnumber := args[3]
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"user\",\"idcardNumber\":\"%s\"}}", idcardnumber)
	queryResults, err := getUsernameForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	if queryResults != nil {
		s := "this idcardnumber already relevanted to another user !!!!"
		ss := []byte(s)
		return shim.Success(ss)
	}
	//查询此手机号是否被注册过
	mobliephone := args[4]
	queryString1 := fmt.Sprintf("{\"selector\":{\"docType\":\"user\",\"mobilePhone\":\"%s\"}}", mobliephone)
	queryResults1, err := getUsernameForQueryString(stub, queryString1)
	if err != nil {
		return shim.Error(err.Error())
	}
	if queryResults1 != nil {
		s := "this mobilephone already relevanted to another user !!!!"
		ss := []byte(s)
		return shim.Success(ss)
	}
	aeskey := string(Krand(16, 3))
	rsapublic := args[5]

	user := &User{objectType, username, usertype, password, idcardnumber, mobliephone, aeskey, rsapublic}
	userJSONasBytes, err := json.Marshal(user)
	err = stub.PutState(username, userJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	s := "1"
	ss := []byte(s)
	return shim.Success(ss)
}

//患者在患者系统登陆
func (t *merChaincode) patientLogin(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	username := args[0]
	password := args[1]
	mobilephone := args[2]
	//查询系统中有无该帐号
	userAsBytes, err := stub.GetState(username)
	if err != nil {
		return shim.Error(err.Error())
	}

	if userAsBytes == nil {
		s := "user is not  exist !!!!"
		ss := []byte(s)
		return shim.Success(ss)
	}

	userTemp := User{}
	json.Unmarshal(userAsBytes, &userTemp)

	if (userTemp.Password != password) || (userTemp.MobilePhone != mobilephone) {
		s := "Login failed  !!!! Please check your information"
		ss := []byte(s)
		return shim.Success(ss)
	}
	s := userTemp.IDcardNumber
	ss := []byte(s)
	return shim.Success(ss)
}

//医生在系统登陆
func (t *merChaincode) doctorLogin(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	username := args[0]
	password := args[1]
	//查询系统中有无该帐号
	userAsBytes, err := stub.GetState(username)
	if err != nil {
		return shim.Error(err.Error())
	}

	if userAsBytes == nil {
		s := "user is not  exist !!!!"
		ss := []byte(s)
		return shim.Success(ss)
	}

	userTemp := User{}
	json.Unmarshal(userAsBytes, &userTemp)

	if userTemp.Password != password {
		s := "Login failed  !!!! Please check your information"
		ss := []byte(s)
		return shim.Success(ss)
	}
	s := userTemp.IDcardNumber
	ss := []byte(s)
	return shim.Success(ss)
}

//新建患者信息
func (t *merChaincode) createPatient(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//获取aes密钥
	username := args[0]
	userAsBytes, err := stub.GetState(username)
	if err != nil {
		return shim.Error(err.Error())
	}
	userTemp := User{}
	json.Unmarshal(userAsBytes, &userTemp)
	aeskey := userTemp.AESKey

	patientId := args[1]
	objectType := "patient"
	timeEn := args[2]
	time := AesDecrypt(timeEn, aeskey)
	nameEn := args[3]
	name := AesDecrypt(nameEn, aeskey)
	gender := AesDecrypt(args[4], aeskey)
	birthday := AesDecrypt(args[5], aeskey)
	nation := AesDecrypt(args[6], aeskey)
	homeAddress := AesDecrypt(args[7], aeskey)
	marriageCondition := AesDecrypt(args[8], aeskey)
	medicalContents := [100]Complaint{}
	reserveInfo := [100]Reservation{}
	doctorIDcollect := [100]string{}
	medicalContentNum := 0
	reserveInfoNum := 0
	doctorIDNum := 0

	//创建患者基本信息对象
	baseinfo := &PatientBaseInfo{patientId, time, name, gender, birthday, nation, homeAddress, marriageCondition}
	//创建患者对象
	patient := &Patient{objectType, patientId, *baseinfo, medicalContents, reserveInfo, doctorIDcollect, medicalContentNum, reserveInfoNum, doctorIDNum}

	patientJSONasBytes, err := json.Marshal(patient)

	//写入账本
	err = stub.PutState(patientId, patientJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	// s := "1"
	// ss := []byte(s)
	return shim.Success(patientJSONasBytes)
}

func getUsernameForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		buffer.WriteString(queryResponse.Key)
	}
	return buffer.Bytes(), nil
}

//查询患者基本信息
func (t *merChaincode) queryPaientBaseinfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	patientId := args[0]
	patientAsBytes, err := stub.GetState(patientId)
	if err != nil {
		return shim.Error(err.Error())
	}
	//获取aes密钥
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"user\",\"idcardNumber\":\"%s\"}}", patientId)
	queryResults, err := getUsernameForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	username := string(queryResults)
	userAsBytes, err := stub.GetState(username)
	if err != nil {
		return shim.Error(err.Error())
	}
	userTemp := User{}
	json.Unmarshal(userAsBytes, &userTemp)
	aeskey := userTemp.AESKey

	patientTemp := Patient{}
	json.Unmarshal(patientAsBytes, &patientTemp)
	baseInfo := patientTemp.BaseInfo

	idcardnumberEn := AesEncrypt(baseInfo.IdCardNumber, aeskey)
	createtimeEn := AesEncrypt(baseInfo.CreateTime, aeskey)
	patientnameEn := AesEncrypt(baseInfo.PatientName, aeskey)
	genderEn := AesEncrypt(baseInfo.Gender, aeskey)
	birthdayEn := AesEncrypt(baseInfo.Birthday, aeskey)
	nationEn := AesEncrypt(baseInfo.Nation, aeskey)
	homeaddressEn := AesEncrypt(baseInfo.HomeAddress, aeskey)
	marriageonditionEn := AesEncrypt(baseInfo.MarriageCondition, aeskey)
	baseinfoEn := &PatientBaseInfo{idcardnumberEn, createtimeEn, patientnameEn, genderEn, birthdayEn, nationEn, homeaddressEn, marriageonditionEn}
	baseInfoEnAsBytes, err := json.Marshal(baseinfoEn)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(baseInfoEnAsBytes)
}

//患者或者医生 更改个人基本信息
func (t *merChaincode) changePaientBaseinfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	patientId := args[0]

	patientAsBytes, err := stub.GetState(patientId)
	if err != nil {
		return shim.Error(err.Error())
	}
	//获取aes密钥
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"user\",\"idcardNumber\":\"%s\"}}", patientId)
	queryResults, err := getUsernameForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	username := string(queryResults)
	userAsBytes, err := stub.GetState(username)
	if err != nil {
		return shim.Error(err.Error())
	}
	userTemp := User{}
	json.Unmarshal(userAsBytes, &userTemp)
	aeskey := userTemp.AESKey

	patientChanged := Patient{}
	json.Unmarshal(patientAsBytes, &patientChanged)
	homeAddress := AesDecrypt(args[1], aeskey)
	marriageCondition := AesDecrypt(args[2], aeskey)
	patientChanged.BaseInfo.HomeAddress = homeAddress
	patientChanged.BaseInfo.MarriageCondition = marriageCondition

	patientJsonAsBytes, _ := json.Marshal(patientChanged)
	err = stub.PutState(patientId, patientJsonAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	s := "1"
	ss := []byte(s)
	return shim.Success(ss)
}

//患者添加预约信息
func (t *merChaincode) addReserveInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	patientId := args[0]
	patientAsBytes, err := stub.GetState(patientId)
	if err != nil {
		return shim.Error(err.Error())
	}
	//获取aes密钥
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"user\",\"idcardNumber\":\"%s\"}}", patientId)
	queryResults, err := getUsernameForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	username := string(queryResults)
	userAsBytes, err := stub.GetState(username)
	if err != nil {
		return shim.Error(err.Error())
	}
	userTemp := User{}
	json.Unmarshal(userAsBytes, &userTemp)
	aeskey := userTemp.AESKey

	patientChanged := Patient{}
	err = json.Unmarshal(patientAsBytes, &patientChanged)
	if err != nil {
		return shim.Error(err.Error())
	}

	num := patientChanged.ReserveInfoNum

	//创建患者的预约信息
	reserveid := AesDecrypt(args[1], aeskey)
	hospitalname := AesDecrypt(args[2], aeskey)
	department := AesDecrypt(args[3], aeskey)
	doctorid := AesDecrypt(args[4], aeskey)
	reserverdate := AesDecrypt(args[5], aeskey)
	reservertime := AesDecrypt(args[6], aeskey)
	reserverstate := AesDecrypt(args[7], aeskey)
	reservation := &Reservation{reserveid, hospitalname, department, doctorid, reserverdate, reservertime, reserverstate}
	//更改patient的相关字段
	patientChanged.ReserveInfo[num] = *reservation
	patientChanged.DoctorIDcollect[num] = args[4]
	num++
	patientChanged.ReserveInfoNum = num
	patientChanged.DoctorIDNum = num
	patientJsonAsBytes, _ := json.Marshal(patientChanged)
	err = stub.PutState(patientId, patientJsonAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	// 更改doctor的相关字段
	doctorAsBytes, err := stub.GetState(doctorid)
	if err != nil {
		return shim.Error(err.Error())
	}

	doctorChanged := Doctor{}
	json.Unmarshal(doctorAsBytes, &doctorChanged)
	number := doctorChanged.PatientIDNum
	doctorChanged.PatientIDcollect[number] = patientId
	number++
	doctorChanged.PatientIDNum = number
	doctorJsonAsBytes, _ := json.Marshal(doctorChanged)
	err = stub.PutState(doctorid, doctorJsonAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	s := "1"
	ss := []byte(s)
	return shim.Success(ss)
}

//查询患者的所有预约数量
func (t *merChaincode) queryReserverInfoNum(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	patientId := args[0]
	patientAsBytes, err := stub.GetState(patientId)
	if err != nil {
		return shim.Error(err.Error())
	}
	//获取aes密钥
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"user\",\"idcardNumber\":\"%s\"}}", patientId)
	queryResults, err := getUsernameForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	username := string(queryResults)
	userAsBytes, err := stub.GetState(username)
	if err != nil {
		return shim.Error(err.Error())
	}
	userTemp := User{}
	json.Unmarshal(userAsBytes, &userTemp)
	aeskey := userTemp.AESKey

	patientTemp := Patient{}
	json.Unmarshal(patientAsBytes, &patientTemp)
	numAsBytes, _ := json.Marshal(patientTemp.ReserveInfoNum)
	encryptCode := AesEncrypt(string(numAsBytes), aeskey)
	return shim.Success([]byte(encryptCode))

}

//按预约ID查询该患者的预约信息
func (t *merChaincode) queryReserverInfoByID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 查询信息有两个参数 即身份证号 病历ID号
	patientId := args[0]
	reserveid := args[1]
	//查询账本中有无该患者
	patientAsBytes, err := stub.GetState(patientId)
	if err != nil {
		return shim.Error(err.Error())
	}

	//获取aes密钥
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"user\",\"idcardNumber\":\"%s\"}}", patientId)
	queryResults, err := getUsernameForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	username := string(queryResults)
	userAsBytes, err := stub.GetState(username)
	if err != nil {
		return shim.Error(err.Error())
	}
	userTemp := User{}
	json.Unmarshal(userAsBytes, &userTemp)
	aeskey := userTemp.AESKey

	patientTemp := Patient{}
	json.Unmarshal(patientAsBytes, &patientTemp)

	reserveinfo := patientTemp.ReserveInfo
	num := patientTemp.ReserveInfoNum

	var res int
	var i int
	for i := 0; i < num; i++ {
		if reserveinfo[i].ReserveID == reserveid {
			res = i
			break
		}
	}
	if i == num {
		s := "this reserve is not exist"
		ss := []byte(s)
		return shim.Success(ss)
	}
	reserve := reserveinfo[res]
	reserveidEn := AesEncrypt(reserve.ReserveID, aeskey)
	hospitalnameEn := AesEncrypt(reserve.HospitalName, aeskey)
	departmentEn := AesEncrypt(reserve.Department, aeskey)
	doctoridEn := AesEncrypt(reserve.DoctorId, aeskey)
	reservedateEn := AesEncrypt(reserve.ReserverDate, aeskey)
	reservetimeEn := AesEncrypt(reserve.ReserverTime, aeskey)
	reservestateEn := AesEncrypt(reserve.ReserverState, aeskey)
	reservationEn := &Reservation{reserveidEn, hospitalnameEn, departmentEn, doctoridEn, reservedateEn, reservetimeEn, reservestateEn}
	reservationEnJsonAsBytes, err := json.Marshal(reservationEn)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(reservationEnJsonAsBytes)
}

//按预约ID更改预约状态
func (t *merChaincode) changeReserverState(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 查询信息有两个参数 即身份证号 病历ID号
	patientId := args[0]
	patientAsBytes, err := stub.GetState(patientId)
	if err != nil {
		return shim.Error(err.Error())
	}
	//获取aes密钥
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"user\",\"idcardNumber\":\"%s\"}}", patientId)
	queryResults, err := getUsernameForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	username := string(queryResults)
	userAsBytes, err := stub.GetState(username)
	if err != nil {
		return shim.Error(err.Error())
	}
	userTemp := User{}
	json.Unmarshal(userAsBytes, &userTemp)
	aeskey := userTemp.AESKey

	reserveid := AesDecrypt(args[1], aeskey)
	patientTemp := Patient{}
	json.Unmarshal(patientAsBytes, &patientTemp)
	reserveinfo := patientTemp.ReserveInfo
	num := patientTemp.ReserveInfoNum
	var res int
	var i int
	for i = 0; i < num; i++ {
		if reserveinfo[i].ReserveID == reserveid {
			res = i
			break
		}
	}
	if i == num {
		s := "this reserve is not exist"
		ss := []byte(s)
		return shim.Success(ss)
	}

	//更改患者字段
	patientTemp.DoctorIDcollect[res] = ""

	// 更改doctor的相关字段
	doctorid := reserveinfo[res].DoctorId
	doctorAsBytes, err := stub.GetState(doctorid)
	if err != nil {
		return shim.Error(err.Error())
	}
	doctorChanged := Doctor{}
	json.Unmarshal(doctorAsBytes, &doctorChanged)
	number := doctorChanged.PatientIDNum
	var j int
	var result int
	for j = 0; j < number; j++ {
		if doctorChanged.PatientIDcollect[j] == patientId {
			result = j
			break
		}
	}
	doctorChanged.PatientIDcollect[result] = ""
	reserveinfo[res] = Reservation{}
	doctorJsonAsBytes, _ := json.Marshal(doctorChanged)
	err = stub.PutState(doctorid, doctorJsonAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	patientTemp.ReserveInfo = reserveinfo
	patientJsonAsBytes, _ := json.Marshal(patientTemp)
	err = stub.PutState(patientId, patientJsonAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	s := "1"
	ss := []byte(s)
	return shim.Success(ss)
}

//删除患者预约信息 （测试用途）
func (t *merChaincode) deleteReserverInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	patientId := args[0]
	//查询账本中有无该患者
	patientAsBytes, err := stub.GetState(patientId)
	if err != nil {
		return shim.Error(err.Error())
	}

	if patientAsBytes == nil {
		return shim.Error("该患者不在系统中")
	}

	patientChanged := Patient{}
	err = json.Unmarshal(patientAsBytes, &patientChanged)
	if err != nil {
		return shim.Error(err.Error())
	}

	patientChanged.ReserveInfo = [100]Reservation{}
	patientChanged.ReserveInfoNum = 0

	patientJsonAsBytes, _ := json.Marshal(patientChanged)
	err = stub.PutState(patientId, patientJsonAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

//检测医生是否有权限对该患者的信息进行查询或者修改
func (t *merChaincode) checkPermission(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	doctorId := args[0]
	patientId := args[1]
	doctorAsBytes, err := stub.GetState(doctorId)
	if err != nil {
		return shim.Error(err.Error())
	}
	doctorTemp := Doctor{}
	json.Unmarshal(doctorAsBytes, &doctorTemp)
	num := doctorTemp.PatientIDNum
	patientidcollect := doctorTemp.PatientIDcollect

	//检测医生是否拥有该患者的权限
	var i int
	for i = 0; i < num; i++ {
		if patientidcollect[i] == patientId {
			break
		}
	}
	if i == num {
		s := "0"
		ss := []byte(s)
		return shim.Success(ss)
	}
	//获取aes密钥
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"user\",\"idcardNumber\":\"%s\"}}", patientId)
	queryResults, err := getUsernameForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	username := string(queryResults)
	userAsBytes, err := stub.GetState(username)
	if err != nil {
		return shim.Error(err.Error())
	}
	userTemp := User{}
	json.Unmarshal(userAsBytes, &userTemp)
	aeskey := userTemp.AESKey

	return shim.Success([]byte(aeskey))
}

//医生添加患者病历
func (t *merChaincode) addMedicalContent(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	doctorId := args[0]
	patientId := args[1]
	patientAsBytes, err := stub.GetState(patientId)
	if err != nil {
		return shim.Error(err.Error())
	}
	//获取aes密钥
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"user\",\"idcardNumber\":\"%s\"}}", patientId)
	queryResults, err := getUsernameForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	username := string(queryResults)
	userAsBytes, err := stub.GetState(username)
	if err != nil {
		return shim.Error(err.Error())
	}
	userTemp := User{}
	json.Unmarshal(userAsBytes, &userTemp)
	aeskey := userTemp.AESKey

	patientChanged := Patient{}
	err = json.Unmarshal(patientAsBytes, &patientChanged)
	if err != nil {
		return shim.Error(err.Error())
	}
	number := patientChanged.MedicalContentNum
	//创建主诉信息
	createdoctorid := doctorId
	complaintid := AesDecrypt(args[2], aeskey)
	department := AesDecrypt(args[3], aeskey)
	time := AesDecrypt(args[4], aeskey)
	medicalType := AesDecrypt(args[5], aeskey)
	symptoms := AesDecrypt(args[6], aeskey)
	conclusion := AesDecrypt(args[7], aeskey)
	presenter := AesDecrypt(args[8], aeskey)
	diseasesOnceSuffered := AesDecrypt(args[9], aeskey)
	medicinecontentNum := 0
	medicineidcontent := [100]string{}
	complaint := &Complaint{complaintid, createdoctorid, department, time, medicalType, symptoms, conclusion, presenter, diseasesOnceSuffered, medicinecontentNum, medicineidcontent}

	patientChanged.MedicalContents[number] = *complaint
	number++
	patientChanged.MedicalContentNum = number
	patientJsonAsBytes, _ := json.Marshal(patientChanged)
	err = stub.PutState(patientId, patientJsonAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	s := "1"
	ss := []byte(s)
	return shim.Success(ss)
}

//删除患者病历  （测试用途）
func (t *merChaincode) deleteMedicalContent(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	patientId := args[0]
	//查询账本中有无该患者
	patientAsBytes, err := stub.GetState(patientId)
	if err != nil {
		return shim.Error(err.Error())
	}

	if patientAsBytes == nil {
		return shim.Error("该患者不在系统中")
	}

	patientChanged := Patient{}
	err = json.Unmarshal(patientAsBytes, &patientChanged)
	if err != nil {
		return shim.Error(err.Error())
	}

	patientChanged.MedicalContents = [100]Complaint{}
	patientChanged.MedicalContentNum = 0

	patientJsonAsBytes, _ := json.Marshal(patientChanged)
	err = stub.PutState(patientId, patientJsonAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

//查询患者的所有病历数量
func (t *merChaincode) queryMedicalNum(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	patientId := args[0]
	//查询账本中有无该患者
	patientAsBytes, err := stub.GetState(patientId)
	if err != nil {
		return shim.Error(err.Error())
	}

	patientTemp := Patient{}
	json.Unmarshal(patientAsBytes, &patientTemp)
	numAsBytes, _ := json.Marshal(patientTemp.MedicalContentNum)
	return shim.Success(numAsBytes)
}

//按病历ID查询该患者的病历
func (t *merChaincode) queryMedicalByID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 查询信息有两个参数 即身份证号 病历ID号
	patientId := args[0]
	complaintid := args[1]
	//查询账本中有无该患者
	patientAsBytes, err := stub.GetState(patientId)
	if err != nil {
		return shim.Error(err.Error())
	}
	//获取aes密钥
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"user\",\"idcardNumber\":\"%s\"}}", patientId)
	queryResults, err := getUsernameForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	username := string(queryResults)
	userAsBytes, err := stub.GetState(username)
	if err != nil {
		return shim.Error(err.Error())
	}
	userTemp := User{}
	json.Unmarshal(userAsBytes, &userTemp)
	aeskey := userTemp.AESKey

	//分离病历信息数组
	patientTemp := Patient{}
	json.Unmarshal(patientAsBytes, &patientTemp)
	medicalcontent := patientTemp.MedicalContents
	num := patientTemp.MedicalContentNum

	var res int
	var i int
	for i = 0; i < num; i++ {
		if medicalcontent[i].ComplaintId == complaintid {
			res = i
			break
		}
	}
	if i == num {
		s := "this medicalcontent is not exist !!!!"
		ss := []byte(s)
		return shim.Success(ss)
	}

	complaintidEn := AesEncrypt(medicalcontent[res].ComplaintId, aeskey)
	createdoctoridEn := AesEncrypt(medicalcontent[res].CreateDoctorId, aeskey)
	departmentEn := AesEncrypt(medicalcontent[res].Department, aeskey)
	medicalcreatetimeEn := AesEncrypt(medicalcontent[res].MedicalCreateTime, aeskey)
	medicaltypeEn := AesEncrypt(medicalcontent[res].MedicalType, aeskey)
	mainsymptomsEn := AesEncrypt(medicalcontent[res].MainSymptoms, aeskey)
	conclusionEn := AesEncrypt(medicalcontent[res].Conclusion, aeskey)
	presenterEn := AesEncrypt(medicalcontent[res].Presenter, aeskey)
	diseasesoncesufferedEn := AesEncrypt(medicalcontent[res].DiseasesOnceSuffered, aeskey)
	medicinecontentnum := medicalcontent[res].MedicineContentNum
	medicineidcontent := medicalcontent[res].MedicineIDContent
	complaintEn := &Complaint{complaintidEn, createdoctoridEn, departmentEn, medicalcreatetimeEn, medicaltypeEn, mainsymptomsEn, conclusionEn, presenterEn, diseasesoncesufferedEn, medicinecontentnum, medicineidcontent}
	complaintEnJsonAsBytes, err := json.Marshal(complaintEn)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(complaintEnJsonAsBytes)
}

//医生为该患者的该条病历开药
func (t *merChaincode) doctorPrescribe(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	patientId := args[0]
	complaintid := args[1]
	patientAsBytes, err := stub.GetState(patientId)
	if err != nil {
		return shim.Error(err.Error())
	}
	//获取aes密钥
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"user\",\"idcardNumber\":\"%s\"}}", patientId)
	queryResults, err := getUsernameForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	username := string(queryResults)
	userAsBytes, err := stub.GetState(username)
	if err != nil {
		return shim.Error(err.Error())
	}
	userTemp := User{}
	json.Unmarshal(userAsBytes, &userTemp)
	aeskey := userTemp.AESKey

	//分离病历信息数组
	patientTemp := Patient{}
	json.Unmarshal(patientAsBytes, &patientTemp)
	medicalcontent := patientTemp.MedicalContents
	num := patientTemp.MedicalContentNum
	var res int
	var i int
	for i = 0; i < num; i++ {
		if medicalcontent[i].ComplaintId == complaintid {
			res = i
			break
		}
	}
	if i == num {
		s := "this medicalcontent is not exist !!!!"
		ss := []byte(s)
		return shim.Success(ss)
	}

	medicineid := AesDecrypt(args[2], aeskey)
	medicinename := AesDecrypt(args[3], aeskey)
	specification := AesDecrypt(args[4], aeskey)
	directions := AesDecrypt(args[5], aeskey)
	remark := AesDecrypt(args[6], aeskey)
	medicine := &Medicine{medicineid, medicinename, specification, directions, remark}
	number := medicalcontent[res].MedicineContentNum
	medicalcontent[res].MedicineIDContent[number] = medicineid
	number++
	medicalcontent[res].MedicineContentNum = number
	patientTemp.MedicalContents = medicalcontent
	patientJsonAsBytes, _ := json.Marshal(patientTemp)
	err = stub.PutState(patientId, patientJsonAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	medicineJsonAsBytes, _ := json.Marshal(medicine)
	err = stub.PutState(medicineid, medicineJsonAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	s := "1"
	ss := []byte(s)
	return shim.Success(ss)
}

//按序号查询药品
func (t *merChaincode) queryMedicineByNum(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	patientId := args[0]
	//获取aes密钥
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"user\",\"idcardNumber\":\"%s\"}}", patientId)
	queryResults, err := getUsernameForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	username := string(queryResults)
	userAsBytes, err := stub.GetState(username)
	if err != nil {
		return shim.Error(err.Error())
	}
	userTemp := User{}
	json.Unmarshal(userAsBytes, &userTemp)
	aeskey := userTemp.AESKey

	medicineid := args[1]
	medicineAsBytes, err := stub.GetState(medicineid)
	if err != nil {
		return shim.Error(err.Error())
	}

	//分离病历信息数组
	medicineTemp := Medicine{}
	json.Unmarshal(medicineAsBytes, &medicineTemp)

	medicineidEn := AesEncrypt(medicineTemp.MedicineId, aeskey)
	medicinenameEn := AesEncrypt(medicineTemp.MedicineName, aeskey)
	specificationEn := AesEncrypt(medicineTemp.Specification, aeskey)
	directionsEn := AesEncrypt(medicineTemp.Directions, aeskey)
	remarkEn := AesEncrypt(medicineTemp.Remark, aeskey)
	medicineEn := &Medicine{medicineidEn, medicinenameEn, specificationEn, directionsEn, remarkEn}
	medicineEnJsonAsBytes, err := json.Marshal(medicineEn)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(medicineEnJsonAsBytes)
}

//新建医生
func (t *merChaincode) createDoctor(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	doctorId := args[0]

	//查询账本中有无该医生
	doctorAsBytes, err := stub.GetState(doctorId)
	if err != nil {
		return shim.Error(err.Error())
	}
	if doctorAsBytes != nil {
		s := "this doctor already existed !!!!"
		ss := []byte(s)
		return shim.Success(ss)
	}

	objectType := "doctor"
	name := args[1]
	hospitalname := args[2]
	role := args[3]
	department := args[4]
	price, _ := strconv.ParseFloat(args[5], 32)
	doctordate := args[6]
	doctortime := args[7]
	patientIDcollect := [100]string{}
	patientIDNum := 0

	//创建医生对象
	doctor := &Doctor{objectType, doctorId, name, hospitalname, role, department, price, doctordate, doctortime, patientIDcollect, patientIDNum}

	doctorJSONasBytes, err := json.Marshal(doctor)
	//写入账本
	err = stub.PutState(doctorId, doctorJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

//自动新建医生
func (t *merChaincode) createDoctorAuto(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("create  start")

	//创建医生对象
	doctor00 := &Doctor{ObjectType: "doctor", DoctorId: "00", DoctorName: "张三", HospitalName: "青岛大学附属医院黄岛分院", Role: "主任医师", Department: "内科", Price: 20.0, DoctorDate: "20190801", DoctorTime: "上午", PatientIDcollect: [100]string{}, PatientIDNum: 0}
	doctor00JSONasBytes, err := json.Marshal(doctor00)
	//写入账本
	err = stub.PutState("00", doctor00JSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	doctor01 := &Doctor{ObjectType: "doctor", DoctorId: "01", DoctorName: "里斯", HospitalName: "青岛大学附属医院黄岛分院", Role: "副主任医师", Department: "内科", Price: 18.0, DoctorDate: "20190801", DoctorTime: "下午", PatientIDcollect: [100]string{}, PatientIDNum: 0}
	doctor01JSONasBytes, err := json.Marshal(doctor01)
	//写入账本
	err = stub.PutState("01", doctor01JSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	doctor02 := &Doctor{ObjectType: "doctor", DoctorId: "02", DoctorName: "王而", HospitalName: "青岛大学附属医院黄岛分院", Role: "主任医师", Department: "外科", Price: 20.0, DoctorDate: "20190802", DoctorTime: "上午", PatientIDcollect: [100]string{}, PatientIDNum: 0}
	doctor02JSONasBytes, err := json.Marshal(doctor02)
	//写入账本
	err = stub.PutState("02", doctor02JSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)

}

//根据医生ID查询医生所有信息
func (t *merChaincode) queryDoctorInfoByID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 查询信息有1个参数 即DoctorID
	doctorId := args[0]
	//查询账本中有无该患者
	doctorAsBytes, err := stub.GetState(doctorId)
	if err != nil {
		return shim.Error(err.Error())
	} else if doctorAsBytes == nil {
		s := "this medicalcontent is not exist !!!!"
		ss := []byte(s)
		return shim.Success(ss)
	}

	return shim.Success(doctorAsBytes)
}

//根据医生ID查询医生的科室
func (t *merChaincode) queryDepartmentByDoctorID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 查询信息有1个参数 即DoctorID
	doctorId := args[0]
	//查询账本中有无该患者
	doctorAsBytes, err := stub.GetState(doctorId)
	if err != nil {
		return shim.Error(err.Error())
	} else if doctorAsBytes == nil {
		s := "this doctor is not exist !!!!"
		ss := []byte(s)
		return shim.Success(ss)
	}

	doctorTemp := Doctor{}
	json.Unmarshal(doctorAsBytes, &doctorTemp)
	departmentAsBytes, _ := json.Marshal(doctorTemp.Department)

	return shim.Success(departmentAsBytes)
}

//查询某一医院的某一科室下的所有医生ID
func (t *merChaincode) queryDoctorByHospitalDepartment(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	hospitalname := args[0]
	department := args[1]
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"doctor\",\"hospitalname\":\"%s\",\"department\":\"%s\"}}", hospitalname, department)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

//转科管理
func (t *merChaincode) transferPermission(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	doctorId := args[0]
	patientId := args[1]
	doctorAsBytes, err := stub.GetState(doctorId)
	if err != nil {
		return shim.Error(err.Error())
	}
	doctorTemp := Doctor{}
	json.Unmarshal(doctorAsBytes, &doctorTemp)
	num := doctorTemp.PatientIDNum

	//检测该医生是否有此患者权限
	for i := 0; i < num; i++ {
		if doctorTemp.PatientIDcollect[i] == patientId {
			s := "this doctor already have permission to this patient !!!!"
			ss := []byte(s)
			return shim.Success(ss)
		}
	}

	doctorTemp.PatientIDcollect[num] = patientId
	num++
	doctorTemp.PatientIDNum = num
	doctorJsonAsBytes, _ := json.Marshal(doctorTemp)
	err = stub.PutState(doctorId, doctorJsonAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	s := "1"
	ss := []byte(s)
	return shim.Success(ss)
}

//couchdb 所需要的查询结果函数
func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"DoctorID\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

func (t *merChaincode) getHistoryForPatient(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	patientid := args[0]

	resultsIterator, err := stub.GetHistoryForKey(patientid)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the marble
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	return shim.Success(buffer.Bytes())
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(merChaincode))
	if err != nil {
		fmt.Printf("Error starting Mer chaincode - %s", err)
	}
}
