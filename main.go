package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	// 1. 커맨드라인 플래그(옵션) 정의
	// -t: 기준 시간 (생략 가능, 기본값: 현재 시간)
	// -d: 더하거나 뺄 시간 (생략 가능, 기본값: 0)
	// -u: 출력 단위 (생략 가능, 기본값: nano)
	timeStr := flag.String("t", "", "기준 시간 (DateTime 형식, 예: '2025-10-18 15:04:05'). 생략 시 현재 시간")
	durationStr := flag.String("d", "0", "계산할 시간 ('-5m', '1h', '300'). 단위 없는 숫자는 초(second)로 간주")
	outputUnit := flag.String("u", "nano", "출력 단위: nano, micro, milli, sec")

	// 프로그램에 전달된 옵션 파싱
	flag.Parse()

	// 2. 기준 시간(baseTime) 결정
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

	// 3. Duration 파싱
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

	// 4. 최종 시간 계산
	finalTime := baseTime.Add(duration)

	// 5. 최종 결과를 지정된 단위로 출력
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
