package ThirdPush

import (
	"fmt"
	"net/smtp"
	"strings"
)

/*
模块说明： 第三方推送接口，邮件推送，短信推送
*/

//PushEmail 邮件推送接口
func PushEmail(toPerson, gatewayName, gatewayID string) {
	auth := smtp.PlainAuth("", "690905084@qq.com", "lqdoamsyipgdbfdc", "smtp.qq.com")
	to := []string{toPerson}
	nickname := "测试"
	user := "690905084@qq.com"
	subject := "网关掉线通知"
	contentType := "Content-Type: text/plain; charset=UTF-8"
	body := fmt.Sprintf("网关名: %s,网关ID:%s ,该设备掉线了，请及时处理!", gatewayName, gatewayID)
	msg := []byte("To: " + strings.Join(to, ",") + "\r\nFrom: " + nickname +
		"<" + user + ">\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + body)
	err := smtp.SendMail("smtp.qq.com:25", auth, user, to, msg)
	if err != nil {
		fmt.Printf("send mail error: %v", err)
	}
	fmt.Println("发送成功")
}
