package mqtt

import (
	"crypto/tls"
	"encoding/json"
	"path/filepath"

	//"encoding/json"
	"fmt"

	"os"
	//"strconv"
	//"sync"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)







//MQTT
var client MQTT.Client = nil


//mqtt配置
type SystemSettingJson struct{
	Server   string `json:"server"`
	Username string `json:"username"`
	Password string `json:"password"`
	//Clientid string `json:"clientid"`

}

func LoadConfig() (SystemSettingJson, bool) {
	var conf SystemSettingJson

	homeTem,err :=os.Getwd()

	if(err!=nil){
		return  conf,false
	}
	ss := homeTem+"/config.json"

	//fmt.Println(CAAInfoPath)

	_dir := filepath.Dir(ss)
	if(!Exists(_dir)){
		//文件夹 不存在
		fmt.Printf("no dir![%v]\n", _dir)
		// 创建文件夹
		err := os.Mkdir(_dir, os.ModePerm)
		if err != nil {

			return  conf,false
		} else {
			fmt.Printf("mkdir Tools success!\n")
		}

	}



	if(Exists(ss)){
		//存在
		filePtr, err := os.Open(ss)
		if err != nil {
			fmt.Println("Open file failed [Err:%s]", err.Error())
			return conf, false
		}
		defer filePtr.Close()
		// 创建json解码器
		decoder := json.NewDecoder(filePtr)
		err = decoder.Decode(&conf)
		if err != nil {
			fmt.Println("[LoadCAAConfig]Decoder failed", err.Error())
			return conf, false

		} else {

			return conf, true
		}
	}else {

		return conf, false
	}

}



////////////////////
//Mqtt初始化
////////////////////
func MqInit() {



	for true{
		for client == nil || client.IsConnected() == false {
			time.Sleep(time.Second*5)
			var server string=""
			var username string=""
			var password string=""
			if jsonData,res:=LoadConfig();res==true{
				server = jsonData.Server
				username = jsonData.Username
				password = jsonData.Password
			}else {
				server = "10.0.0.142:1883"
				username = ""
				password = ""
			}




			connOpts := MQTT.NewClientOptions().AddBroker(server).SetClientID("ringdoor2mqtt").SetCleanSession(true)
			if username != "" {
				connOpts.SetUsername(username)
				if password != "" {
					connOpts.SetPassword(password)
				}
			}
			tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
			connOpts.SetTLSConfig(tlsConfig)


			client = MQTT.NewClient(connOpts)
			if token := client.Connect(); token.Wait() && token.Error() != nil {

				fmt.Println("[ERROR] mqtt connect fail!"+token.Error().Error())

				//return false
				continue
			}

			fmt.Println("mqtt connect success")

		}
	}



}

func Mqpublic(topic string, qos int, retained bool, message string) bool {
	if client == nil || client.IsConnected() == false {
		//LogPublish("Send fail!Uable to Reconnect the mqtt server!\n", common.LogError)
		//发送失败

		return false

	}

	//sendTem,_ :=json.Marshal(message)
	if token := client.Publish(topic, byte(qos), retained, message); token.Wait() && token.Error() != nil {

		fmt.Println(token.Error())
		return false
	}

	return true
}

////////////////////////////
//文件/文件夹是否存在
/////////////////////////
func Exists(path string) bool {
	_, err := os.Stat(path)    //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}




