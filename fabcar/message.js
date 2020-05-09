const Core = require('@alicloud/pop-core');

function sendMessage(number,code)
{
	var client = new Core({
		accessKeyId: '00000',
		accessKeySecret: '000000',
		endpoint: 'https://dysmsapi.aliyuncs.com',
		apiVersion: '2017-05-25'
	  });
	  
	  var params = {
		"RegionId": "cn-hangzhou",
		"PhoneNumbers": number,
		"SignName": "Ecase",
		"TemplateCode": "SMS_172887364",
		"TemplateParam":code
	  }
	  console.log(params)
	  var requestOption = {
		method: 'POST'
	  };
	  
	  client.request('SendSms', params, requestOption).then((result) => {
		console.log(JSON.stringify(result));
		return JSON.stringify(result);
	  }, (ex) => {
		console.log(ex);
	  })

}

module.exports.sendMessage = sendMessage;