package sms

import (
	"encoding/json"
	"errors"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"go-msg/model"
)

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func CreateAliClient(accessKeyId *string, accessKeySecret *string) (_result *dysmsapi20170525.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId: accessKeyId,
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	_result = &dysmsapi20170525.Client{}
	_result, _err = dysmsapi20170525.NewClient(config)
	return _result, _err
}


func SendAliSms(phone string, params interface{}, template model.SmsTemplate) error {
	client, _err := CreateAliClient(tea.String(template.AccessKeyId), tea.String(template.AccessKeySecret))
	if _err != nil {
		return _err
	}
	paramsJson, err := json.Marshal(params)
	if err != nil {
		return errors.New("短信发送失败，参数异常")
	}
	paramsJsonStr := string(paramsJson)
	signName := template.SignName
	templateCode := template.TemplateCode
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		SignName:      &signName,
		TemplateCode:  &templateCode,
		TemplateParam: &paramsJsonStr,
		PhoneNumbers:  &phone,
	}
	_, _err = client.SendSms(sendSmsRequest)
	if _err != nil {
		return _err
	}
	return _err
}



type SendSmsCommonReq struct {
	Mobile string `alias:"手机号码" json:"mobile" validate:"number,startswith=1,len=11"`
	Op     string `alias:"操作" json:"op" validate:"required"`
	OrderId  string `alias:"参数" json:"orderId,omitempty"`
	Consignee string `alias:"乘客姓名" json:"consignee,omitempty" `
	Phone string `alias:"手机号码" json:"phone,omitempty"`
	Date string `alias:"时间" json:"date,omitempty"`
	Name string `alias:"司机名称" json:"name,omitempty"`
	Transfer string `alias:"行程名" json:"transfer,omitempty"`
}


func Send(){
	// template := model.SmsTemplate{}
	// tx := .Model(model.SmsTemplate{}).First(&template, "code=?", s.SendCommonOrderDto.Op)
	// if tx.Error != nil {
	// 	ctx.Log.Error(tx.Error)
	// 	return errcode.ExecuteErr("短信发送失败，模板获取异常")
	// }
	// 2. 组织发送的数据body
    // smsParams := util.Struct2Map(s.SendCommonOrderDto)
	// delete(smsParams,"mobile")
    // delete(smsParams,"op")
	// logger.GetInstance().Info("**** 发送短信的模板参数：",packages.ToRecordJson(smsParams))
	// 3. 发送短信
	// err := sms_provider.SendAliSms(s.SendCommonOrderDto.Mobile, smsParams, template)
	// if err != nil {
	// 	ctx.Log.Error(err)
	// 	return errcode.ParamsErr("短信发送失败")
	// }
	// return nil
}
