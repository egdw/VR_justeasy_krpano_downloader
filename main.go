package main

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var locations []string = []string{"b", "d", "f", "l", "r", "u"}

type Scene struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Son  []struct {
		Name         string `json:"name"`
		Groupid      int    `json:"groupid"`
		Sceneid      string `json:"sceneid"`
		Pano2Sceneid string `json:"pano2sceneid"`
		Preview      string `json:"preview"`
		Hide         int    `json:"hide"`
	} `json:"son"`
}

//得到所有的场景
func getAllScene(url string) []Scene {
	client := http.Client{}
	resp, err := client.Get(url)
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}
	readAll, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	//查找json_list位置
	jsonStartIndex := strings.Index(string(readAll), "}],jsonList:")
	jsonStart := string(readAll)[jsonStartIndex+11:]
	//查找json结尾位置
	jsonEndIndex := strings.Index(jsonStart, "}]}]};")
	var Scene []Scene
	json.Unmarshal([]byte(string(readAll)[jsonStartIndex+12:jsonStartIndex+15+jsonEndIndex]), &Scene)
	return Scene
}

//判断文件夹是否存在
func HasDir(path string) (bool, error) {
	_, _err := os.Stat(path)
	if _err == nil {
		return true, nil
	}
	if os.IsNotExist(_err) {
		return false, nil
	}
	return false, _err
}

//创建文件夹
func CreateDir(path string) {
	_exist, _err := HasDir(path)
	if _err != nil {
		fmt.Printf("获取文件夹异常 -> %v\n", _err)
		return
	}
	if _exist {
		//fmt.Println("文件夹已存在！")
	} else {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			fmt.Printf("创建目录异常 -> %v\n", err)
		}
	}
}

//用于解析图片数据
func ParserData(name string, panoUrl string) {
	date := time.Now().Format("20060102")
	CreateDir("./output/" + name + date)
	for _, location := range locations {
		fileName := "/l1_" + location + "_1_1.jpg"
		downloadUrl := "https://vrimg.justeasy.cn/" + strings.ReplaceAll(panoUrl, "thumb.jpg", location+fileName)
		if err := downloadFile(name, downloadUrl, fileName); err != nil {
			fmt.Println(name+":"+fileName+"下载失败：", err)
		}
		valid, err := isValidJPEG("./output/" + name + date + fileName)
		if !valid {
			fmt.Println(name+":"+fileName+"，文件验证失败，将尝试备用地址下载，", err)
			downloadUrl = "https://vrpic.justeasy.cn/" + strings.ReplaceAll(panoUrl, "thumb.jpg", location+fileName)
			if err = downloadFile(name, downloadUrl, fileName); err != nil {
				fmt.Println(name+":"+fileName+"下载失败：", err)
			}
			fmt.Println(name + ":" + fileName + "下载成功，备用地址下载成功。")
		}
		//println(downloadUrl)
	}
	CubeToSphere(".\\output\\" + name + date + "\\")
}

func isValidJPEG(filename string) (bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return false, fmt.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()

	_, format, err := image.Decode(file)
	if err != nil {
		return false, fmt.Errorf("JPEG文件损坏或格式不正确: %v", err)
	}
	if format != "jpeg" {
		return false, fmt.Errorf("文件不是JPEG格式，检测到的格式: %s", format)
	}

	return true, nil
}

func downloadFile(name, url, fileName string) error {
	date := time.Now().Format("20060102")
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile("./output/"+name+date+fileName, data, 0644); err != nil {
		return err
	}
	return nil
}

//将方形图转为全景图
func CubeToSphere(dirUrl string) {
	//systemBit := fmt.Sprint(32 << (^uint(0) >> 63))
	//println(systemBit)
	//commandStr := fmt.Sprintf("module\\krpanotools%s.exe cube2sphere -config=module\\convertdroplets.config \"%s\" \"%s\" \"%s\" \"%s\" \"%s\" \"%s\" ", systemBit,
	//	dirUrl+"l1_b_1_1.jpg", dirUrl+"l1_d_1_1.jpg", dirUrl+"l1_f_1_1.jpg", dirUrl+"l1_l_1_1.jpg", dirUrl+"l1_r_1_1.jpg", dirUrl+"l1_u_1_1.jpg")
	//commandStr := fmt.Sprintf("module\\krpanotools%s.exe", systemBit)
	//println(commandStr)

	command := exec.Command("cmd.exe", "/C", "module\\cube2sphere.bat", dirUrl+"l1_b_1_1.jpg", dirUrl+"l1_d_1_1.jpg", dirUrl+"l1_f_1_1.jpg", dirUrl+"l1_l_1_1.jpg", dirUrl+"l1_r_1_1.jpg", dirUrl+"l1_u_1_1.jpg")
	combinedOutput, err := command.CombinedOutput()
	if err != nil {
		panic(err)
	}
	commandResult := string(combinedOutput)
	//println(commandResult)
	if strings.Index(commandResult, "l1_sphere.jpg") >= 0 {
		exec.Command("start", dirUrl).Start()
	} else {

	}
}

func registerTools() {
	command := exec.Command("cmd.exe", "/C", "module\\activate.bat")
	combinedOutput, err := command.CombinedOutput()
	if err != nil {
		panic(err)
	}
	commandResult := string(combinedOutput)
	if strings.Index(commandResult, "Code registered.") < 0 {
		println("激活失败，生成的全景图会存在水印！")
	}

}

func main() {
	println("vr.justeasy.cn全景图下载工具--恶搞大王")
	registerTools() //激活工具，防止出现水印
	println("==============请粘贴需要下载的全景图网站地址===============")
	println("例如：https://vr.justeasy.cn/view/1656cods54814n65-1658438750.html")
	var url string
	print("请输入全景图地址：")
	fmt.Scanln(&url)
	url = strings.Trim(url, "") //去除空格
	if url == "" {
		println("地址有误！")
		return
	}
	scenes := getAllScene(url)
	if scenes == nil || len(scenes) == 0 {
		println("没有找到场景！")
	}
	if len(scenes) == 1 {
		//说明只有一个场景
		println("共找到以下几个场景：")
		for index, scene := range scenes[0].Son {
			println(strconv.Itoa(index+1), ":", scene.Name)
		}
		println("==============请选择导出场景===============")
		println("0:全部导出，{0-N}:导出某个场景。例如输入1导出第一个场景。")
		outputSceneMethod := 0
		print("请输入导出方式：")
		fmt.Scanln(&outputSceneMethod)
		if outputSceneMethod == 0 {
			for i := 0; i < len(scenes[0].Son); i++ {
				println("-------------------->下载", scenes[0].Son[i].Name, "中<--------------------")
				ParserData(scenes[0].Son[i].Name, scenes[0].Son[i].Preview)
				println("下载", scenes[0].Son[i].Name, "完成。")
			}
		} else {
			if outputSceneMethod > len(scenes[0].Son) {
				println("输入场景编号有误！")
				return
			}
			if outputSceneMethod < 0 {
				println("场景编号不能小于0")
				return
			}
			ParserData(scenes[0].Son[outputSceneMethod-1].Name, scenes[0].Son[outputSceneMethod-1].Preview)
		}
	} else {
		//说明有多个场景
		println("暂不支持识别多个场景！")

	}
	println("请到如下地址进行全景图预览：")
	println("https://egdw.github.io/panorama_display/")
	var pause string
	fmt.Scanln(&pause)
}
