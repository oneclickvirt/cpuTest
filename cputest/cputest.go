package cputest

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/parnurzeal/gorequest"
)

// runSysBenchCommand 执行 sysbench 命令进行测试
func runSysBenchCommand(numThreads, maxTime, version string) (string, error) {
	// version <= 1.0.17
	// sysbench --test=cpu --num-threads=1 --cpu-max-prime=10000 --max-requests=1000000 --max-time=5 run
	// version >= 1.0.18
	// sysbench cpu --threads=1 --cpu-max-prime=10000 --events=1000000 --time=5 run
	var command *exec.Cmd
	if strings.Contains(version, "1.0.18") || strings.Contains(version, "1.0.19") || strings.Contains(version, "1.0.20") {
		command = exec.Command("sysbench", "cpu", "--threads="+numThreads, "--cpu-max-prime=10000", "--events=1000000", "--time="+maxTime, "run")
	} else {
		command = exec.Command("sysbench", "--test=cpu", "--num-threads="+numThreads, "--cpu-max-prime=10000", "--max-requests=1000000", "--max-time="+maxTime, "run")
	}
	output, err := command.CombinedOutput()
	return string(output), err
}

func SysBenchTest(language, testThread string) string {
	var result, singleScore, multiScore string
	comCheck := exec.Command("sysbench", "--version")
	output, err := comCheck.CombinedOutput()
	if err == nil {
		version := string(output)
		if testThread == "single" {
			singleResult, err := runSysBenchCommand("1", "5", version)
			if err == nil {
				tempList := strings.Split(singleResult, "\n")
				for _, line := range tempList {
					if strings.Contains(line, "events per second:") || strings.Contains(line, "total number of events:") {
						temp1 := strings.Split(line, ":")
						if len(temp1) == 2 {
							singleScore = temp1[1]
						}
					}
				}
			}
			if language == "en" {
				result += "1 Thread(s) Test: "
			} else {
				result += "1 线程测试(单核)得分: "
			}
			result += singleScore + "\n"
		} else if testThread == "multi" {
			singleResult, err := runSysBenchCommand("1", "5", version)
			if err == nil {
				tempList := strings.Split(singleResult, "\n")
				for _, line := range tempList {
					if strings.Contains(line, "events per second:") || strings.Contains(line, "total number of events:") {
						temp1 := strings.Split(line, ":")
						if len(temp1) == 2 {
							singleScore = temp1[1]
						}
					}
				}
			}
			if language == "en" {
				result += "1 Thread(s) Test: "
			} else {
				result += "1 线程测试(单核)得分: "
			}
			result += singleScore + "\n"
			multiResult, err := runSysBenchCommand(fmt.Sprintf("%s", runtime.NumCPU()), "5", version)
			if err == nil {
				tempList := strings.Split(multiResult, "\n")
				for _, line := range tempList {
					if strings.Contains(line, "events per second:") || strings.Contains(line, "total number of events:") {
						temp1 := strings.Split(line, ":")
						if len(temp1) == 2 {
							multiScore = temp1[1]
						}
					}
				}
			}
			if language == "en" {
				result += fmt.Sprintf("%d", runtime.NumCPU()) + " Thread(s) Test: "
			} else {
				result += fmt.Sprintf("%d", runtime.NumCPU()) + " 线程测试(多核)得分: "
			}
			result += multiScore + "\n"
		}
	}
	return result
}

// runGeekbenchCommand 执行 geekbench 命令进行测试
func runGeekbenchCommand() (string, error) {
	var command *exec.Cmd
	command = exec.Command("/tmp/geekbench/geekbench", "--upload")
	output, err := command.CombinedOutput()
	return string(output), err
}

// GeekBenchTest 调用 geekbench 执行CPU测试
// 调用 /tmp 下的 /tmp/geekbench 中的 geekbench 文件执行
// https://github.com/masonr/yet-another-bench-script/blob/0ad4c4e85694dbcf0958d8045c2399dbd0f9298c/yabs.sh#L894
func GeekBenchTest(language, testThread string) string {
	var result, singleScore, multiScore, link string
	comCheck := exec.Command("/tmp/geekbench/geekbench", "--version")
	// Geekbench 5.4.5 Tryout Build 503938 (corktown-master-build 6006e737ba)
	output, err := comCheck.CombinedOutput()
	version := string(output)
	if err == nil {
		if strings.Contains(version, "Geekbench 6") {
			// 检测存在 /etc/os-release 文件且含 CentOS Linux 7 时，需要预先下载 GLIBC_2.27 才能使用 geekbench 6
			file, err := os.Open("/etc/os-release")
			defer file.Close()
			if err == nil {
				scanner := bufio.NewScanner(file)
				isCentOS7 := false
				// 逐行读取文件内容
				for scanner.Scan() {
					line := scanner.Text()
					if strings.Contains(line, "CentOS Linux 7") {
						isCentOS7 = true
						break
					}
				}
				if err := scanner.Err(); err == nil {
					// 如果文件中包含 CentOS Linux 7，则打印提示信息
					if isCentOS7 && language == "zh" {
						return "需要预先下载 GLIBC_2.27 才能使用 geekbench 6"
					} else if isCentOS7 && language != "zh" {
						return "You need to pre-download GLIBC_2.27 to use geekbench 6."
					}
				}
			}
		}
		// 解析 geekbench 执行结果
		if strings.Contains(version, "Geekbench") {
			tp, err := runGeekbenchCommand()
			if err == nil {
				tempList := strings.Split(tp, "\n")
				for _, line := range tempList {
					if strings.Contains(line, "https://browser.geekbench.com") && strings.Contains(line, "cpu") {
						link = strings.TrimSpace(line)
						break
					}
				}
			}
		}
		const (
			userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36"
			accept    = "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"
			referer   = "browser.geekbench.com"
		)
		request := gorequest.New().Get(link).
			Retry(3, 6*time.Second, http.StatusBadRequest, http.StatusInternalServerError)
		request.Set("User-Agent", userAgent)
		request.Set("Accept", accept)
		request.Set("Referer", referer)
		response, body, err := request.End()
		if err != nil {
			return ""
		}
		if response.StatusCode != http.StatusOK {
			return ""
		}
		doc, readErr := goquery.NewDocumentFromReader(strings.NewReader(body))
		if readErr != nil {
			return ""
		}
		textContent := doc.Find(".table-wrapper.cpu").Text()
		resList := strings.Split(textContent, "\n")
		for index, l := range resList {
			if strings.Contains(l, "Single-Core") {
				singleScore = resList[index-1]
			} else if strings.Contains(l, "Multi-Core") {
				multiScore = resList[index-1]
			}
		}
	}
	if link != "" && singleScore != "" {
		result += strings.TrimSpace(strings.ReplaceAll(version, "\n", "")) + "\n"
		result += "Single-Core Score: " + singleScore + "\n"
		if multiScore != "" {
			result += "Multi-Core Score: " + multiScore + "\n"
		}
		result += "Link: " + link + "\n"
	}
	return result
}

func WinsatTest(language, testThread string) string {
	var result string
	cmd1 := exec.Command("winsat", "cpu", "-encryption")
	output1, err1 := cmd1.Output()
	if err1 != nil {
		return ""
	} else {
		tempList := strings.Split(string(output1), "\n")
		for _, l := range tempList {
			if strings.Contains(l, "CPU AES256") {
				tempL := strings.Split(l, " ")
				tempText := strings.TrimSpace(tempL[len(tempL)-2])
				if language == "en" {
					result += "CPU AES256 encrypt: "
				} else {
					result += "CPU AES256 加密: "
				}
				result += tempText + "MB/s" + "\n"
			}
		}
	}
	cmd2 := exec.Command("winsat", "cpu", "-compression")
	output2, err2 := cmd2.Output()
	if err2 != nil {
		return ""
	} else {
		tempList := strings.Split(string(output2), "\n")
		for _, l := range tempList {
			if strings.Contains(l, "CPU LZW") {
				tempL := strings.Split(l, " ")
				tempText := strings.TrimSpace(tempL[len(tempL)-2])
				if language == "en" {
					result += "CPU LZW Compression: "
				} else {
					result += "CPU LZW 压缩: "
				}
				result += tempText + "MB/s" + "\n"
			}
		}
	}
	return result
}
