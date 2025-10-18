package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	// -t: 기준 시간 (생략 가능, 기본값: 현재 시간)
	// -d: 더하거나 뺄 시간 (생략 가능, 기본값: 0)
	// -u: 출력 단위 (생략 가능, 기본값: nano)
	timeStr := flag.String("t", "", "- 기준 시간 (DateTime 형식, 예: '2025-10-18 15:04:05'). 생략 시 현재 시간\n- ex: ug -t 2025-10-19 01:24:00\n- '2025-10-19 01:24:00' 시간을 unix time nano 단위로 변환")
	durationStr := flag.String("d", "0", "- 계산할 시간 ('-5m', '1h', '300'). 단위 없는 숫자는 초(second)로 간주\n- ex: ug -d -10m\n- 현재 시간의 10분 전 시간을 unixtime nano 단위로 변환")
	outputUnit := flag.String("u", "nano", "- 출력 단위: nano, micro, milli, sec\n- ex: ug -d 10m -u milli\n- 현재 시간의 10분 후 시간을 unixtime milli 단위로 변환")
	toDateTime := flag.Int64("dt", 0, "- unixtime을 datetime으로 변환\n- ex: ug -dt 1760807714195896000\n- 이 옵션 사용시 'u' 옵션만 사용 가능, 다른 옵션(d, t)은 무시 됨\n- 옵션을 생략하면 자동으로 nano단위로 변환하고 다른 단위의 unixtime을 입력할 경우 반드시 'u' 옵션을 통해 단위를 지정해줘야 함")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "### 아무 옵션 없이 'ug' 실행시 현재 시간 기준으로 unixtime nano 출력 ###\n\n")
		flag.PrintDefaults() // 기존 플래그 옵션 출력
	}

	flag.Parse()

	if *toDateTime != 0 {
		switch *outputUnit {
		case "nano":
			sec := *toDateTime / 1e9
			nsec := *toDateTime % 1e9
			fmt.Println(time.Unix(sec, nsec))
		case "micro":
			fmt.Println(time.UnixMicro(*toDateTime))
		case "milli":
			fmt.Println(time.UnixMilli(*toDateTime))
		case "sec":
			fmt.Println(time.Unix(*toDateTime, 0))
		default:
			fmt.Fprintf(os.Stderr, "오류: 잘못된 출력 단위 '%s'입니다. 'nano', 'micro', 'milli', 'sec' 중 하나를 사용하세요.\n", *outputUnit)
			os.Exit(1)
		}
		return
	}

	var baseTime time.Time
	var err error

	if *timeStr == "" {
		// -t 옵션이 없으면 현재 시간 사용
		baseTime = time.Now()
	} else {
		// -t 옵션이 있으면 해당 문자열을 시간으로 파싱
		baseTime, err = time.Parse(time.DateTime, *timeStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "오류: -t 옵션의 시간 형식이 잘못되었습니다: %v\n", err)
			os.Exit(1)
		}
	}

	var duration time.Duration

	// 먼저 "5m", "-1h30m" 같은 표준 Duration 형식으로 파싱 시도
	duration, err = time.ParseDuration(*durationStr)
	if err != nil {
		// 파싱 실패 시, 단위 없는 숫자(초) 형식으로 파싱 시도
		seconds, err := strconv.ParseInt(*durationStr, 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "오류: -d 옵션의 형식이 잘못되었습니다. '5m', '-1h' 또는 '300' 같은 형식을 사용하세요.\n")
			os.Exit(1)
		}
		duration = time.Duration(seconds) * time.Second
	}

	finalTime := baseTime.Add(duration)

	fmt.Printf("🔹 기준 시간: %s\n", baseTime.Format(time.RFC3339Nano))
	fmt.Printf("🔹 적용 시간: %s\n", duration)
	fmt.Printf("✅ 결과 시간: %s\n", finalTime.Format(time.RFC3339Nano))
	fmt.Print("🔢 출력 값: ")

	switch *outputUnit {
	case "nano":
		fmt.Println(finalTime.UnixNano())
	case "micro":
		fmt.Println(finalTime.UnixMicro())
	case "milli":
		fmt.Println(finalTime.UnixMilli())
	case "sec":
		fmt.Println(finalTime.Unix())
	default:
		fmt.Fprintf(os.Stderr, "오류: 잘못된 출력 단위 '%s'입니다. 'nano', 'micro', 'milli', 'sec' 중 하나를 사용하세요.\n", *outputUnit)
		os.Exit(1)
	}
}
