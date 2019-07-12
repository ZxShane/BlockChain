package main

import (
//	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
//	"strings"
//	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type merChaincode struct {
}


//  患者病历模型设计
type  Patient struct{
	ObjectType string        `json:"objectType"` //field for couchdb
	CreateTime  int64              `json:"createtime"`      //病理创建时间
    IdCardNumber string      `json:"idcardnumber"`     //患者身份证号
	PatientName string                `json:"patientName"` 
	Gender string              `json:"gender"`     //性别m为男
	Birthday  string           `json:"birthday"` 
	Nation string             `json:"nation"` 
	HomeAddress string           `json:"homeAddress"` 
	AllowQuery   int              `json:"alllowQuery"`      //是否有权限查询 0不允许
	AllowAppend   int            `json:"allowAppend"` 
    MedicalContents []Complaint   `json:"medicalContents"`    //病历数组
}


//患者主诉信息
type Complaint struct{
	 Idnumber int           `json:"idnumber"`           //患者该条病历的 ID
	 AllowModify  int          `json:"allowModify"`     //是否能够修改病历
	 MedicalCreateTime  int64     `json:"medicalCreate"`     //该条病历的创建时间
	 MedicalType int64    `json:"medicalType"`           //患者疾病类型
	 MainSymptoms string       `json:"mainSymptoms"`    //患者主诉
	 DetailSymptoms    string    `json:"detailSymptoms"`      //患者症状详细信息   *************************
	 DiseasesOnceSuffered   string   `json:"diseasesOnceSuffered"`    //患者曾患病
	 SystemReview  string   `json:"systemReview"`     //系统回顾
	 PatientPersonalDesctiption  string   `json:"patientPersonalDesctiption"`    //患者个人史***************
}


//患者症状详细描述信息
type DetailSymptomsContent struct{
	OnesetTimeAndPossiblePathogeny   string         `json:"onesetTimeAndPossiblePathogeny"`   //发病时间和可能原因
	MainSymptomsElaborateDescription  string          `json:"mainSymptomsElaborateDescription"`     //主要症状的详细描述
	SimultaneousPhenomenon string     `json:"simultaneousPhenomenon"`         //伴随症状
    OtherDiseases string    `json:"otherDiseases"`     //患者患有的其他疾病
	GeneralConditions string    `json:"generalConditions"`    //发病以来的一般情况
}


//患者个人史信息
type PatientPersonalDescription struct{
	RecentAreaAndDate string          `json:"recentAreaAndDate"`    //患者最近去过的地点以及时间
	HabitsAndCustoms   string           `json:"habitsAndCustoms"`     //患者起居习惯等
	Occupation   string             `json:"occupation"`       //患者的职业和工作环境
	VisitProstitutes  string    `json:"visitProstitutes"`      //患者有无冶游史，以及发病时间
	GrowthAndDevelopmentHistory     string    `json:"growthAndDevelopmentHistory"`     //患者生长发育史
	MaritalStaus  string     `json:"maritalStaus"`      //患者婚姻情况
	Menstruation  string    `json:"menstruation"`      //女性患者月经情况

 }


 //链代码初始化函数
 func (t *merChaincode) Init(stub shim.ChaincodeStubInterface)  pb.Response {
	fmt.Println("Patient Init Success")
	return shim.Success(nil)
}

//Invoke 函数
func (t *merChaincode)  Invoke(stub shim.ChaincodeStubInterface)  pb.Response {
	function, args := stub.GetFunctionAndParameters()
	
	//根据function 参数值调用相应的处理函数
	if  function == "create"{
		return t.create(stub,args)
	} else if function == "isAllowQueryUserContent" {
		return t.isAllowQueryUserContent(stub,args)
	}  else  if  function == "patientRegistration" {
		return  t.patientRegistration(stub, args)
	}
	return shim.Error("没找到对应方法~")
}

//是否可被查询
func (t *merChaincode)  isAllowQueryUserContent(stub shim.ChaincodeStubInterface,args [] string)  pb.Response{
    // 只有一个参数 即身份证号
	patientId := args[0]
	//查询账本中有无该用户
	patientAsBytes, err := stub.GetState(patientId)
	if err != nil {
		return shim.Error(err.Error())
	}else if patientAsBytes == nil {
		return shim.Error("该用户不存在~")
	  }
	  
	  //查询有此人 获取查询权值进行判断
	 patientInstance := Patient{}
	 err = json.Unmarshal(patientAsBytes ,  &patientInstance)
	 if err != nil {
		 return shim.Error(err.Error())
	 }
	  
	 if patientInstance.AllowQuery == 1{
		 return shim.Success(nil) 
	 } else {
         return shim.Error("不可被查询！！")
	 }
     
}


//是否可被修改
func (t *merChaincode)  isAllowAppendUserContent(stub shim.ChaincodeStubInterface,args [] string)  pb.Response{
    // 只有一个参数 即身份证号
	patientId := args[0]
	//查询账本中有无该用户
	patientAsBytes, err := stub.GetState(patientId)
	if err != nil {
		return shim.Error(err.Error())
	}else if patientAsBytes == nil {
		return shim.Error("该用户不存在~")
	  }
	  
	  //查询有此人 获取查询权值进行判断
	 patientInstance := Patient{}
	 err = json.Unmarshal(patientAsBytes ,  &patientInstance)
	 if err != nil {
		 return shim.Error(err.Error())
	 }
	  
	 if patientInstance.AllowAppend== 1{
		 return shim.Success(nil) 
	 } else {
         return shim.Error("不可被追加病历！！")
	 }
     
}



// 解锁用户病历的可修改 可增加权限
func (t *merChaincode)  patientRegistration(stub shim.ChaincodeStubInterface,args [] string)  pb.Response{
	 
	 //三个参数  都为 int 类型 改变 AllowQuery 和 AllowAppend 字段
	 
	 patientId := args[0]
	 //查询账本中有无该用户
	patientAsBytes, err := stub.GetState(patientId)
	if err != nil {
		return shim.Error(err.Error())
	}else if patientAsBytes == nil {
		return shim.Error("该用户不存在~")
	  }

	 //查询有此人 修改权限
	 patientChanged := Patient{}
	 err = json.Unmarshal(patientAsBytes ,  &patientChanged)
	 if err != nil {
		 return shim.Error(err.Error())
	 }

	new_QueryState ,err := strconv.Atoi(args[1] )
	new_AppendState , err := strconv.Atoi(args[2])
	patientChanged.AllowQuery = new_QueryState
	patientChanged.AllowAppend = new_AppendState

	patientJsonAsBytes , _ := json.Marshal(patientChanged)
	err = stub.PutState(patientId , patientJsonAsBytes)
    if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}


//为患者创建新病历
func (t *merChaincode)  create(stub shim.ChaincodeStubInterface,args [] string)  pb.Response {
	fmt.Println("create  start")
	
	
	var err error 
	patientId := args[0]

	//查询账本中有无该用户
	patientIdAsBytes, err := stub.GetState(patientId)
	if err != nil {
		return shim.Error(err.Error())
	}
	
	if  patientIdAsBytes != nil {
		return shim.Error("该用户的病历已经存在，无需创建新病历！！")
	  }
	
	time,err := strconv.ParseInt(args[0] ,10,64)
	objectType := "patient"
	name := args[1]
	gender := args[2]
	birthday := args[3]
	nation := args[4]
	homeAddress := args[5]

	allowquery,err := strconv.Atoi(args[7])
	allowappend,err := strconv.Atoi(args[8])
	medicalContents := [] Complaint{}

	//patient := &Patient{ objectType,patientId,name,gender,birthday,nation,homeAddress}
	patient := &Patient{ objectType,time,patientId,name,gender,birthday,nation,homeAddress,allowquery,allowappend, medicalContents}
	
	
	patientJSONasBytes, err := json.Marshal(patient)
	//写入账本
	err = stub.PutState(patientId,patientJSONasBytes)
	if err != nil{
		return shim.Error(err.Error())
	}

	return shim.Success(nil)


}


//添加患者病历
func (t *merChaincode)  addNewMedicalContent(stub shim.ChaincodeStubInterface,args [] string)  pb.Response {
	patientId := args[0]

	//查询账本中有无该用户
	patientIdAsBytes, err := stub.GetState(patientId)
	if err != nil {
		return shim.Error(err.Error())
	}
	
	if  patientIdAsBytes == nil {
		return shim.Error("该用户还没有病历，请先创建病历")
	  }
	 
	  complaintID , err := strconv.Atoi(args[1])
	  allow_modify , err := strconv.Atoi(args[2])
	  time,err := strconv.ParseInt(args[3] ,10,64)
	  medical_type ,err := strconv.ParseInt(args[4] ,10,64)
	  symptoms := args[5]
	  detailSymptoms := args[6]
	  diseasesOnceSuffered := args[7]
	  systemReview := args[8]
	  patientPersonalDesctiption := args[9]

	  complaint := &Complaint{complaintID, allow_modify,time,medical_type,symptoms,detailSymptoms,diseasesOnceSuffered,systemReview,patientPersonalDesctiption}

	  complaintJSONasBytes, err := json.Marshal(complaint)
	//写入账本
	err = stub.PutState(args[1],complaintJSONasBytes)
	if err != nil{
		return shim.Error(err.Error())
	}

	return shim.Success(nil)

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


