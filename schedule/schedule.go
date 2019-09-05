package schedule

import (
	"log"
	"loranet20181205/database"
	"loranet20181205/webdata"
	"time"
)

// Schedule1 check every 15 minutes
func Schedule1(interval int64, InitDelay int64) {
	/* initial delay */
	time.Sleep(time.Duration(InitDelay) * time.Second)
	for {

		/* put the schedule code here */
		//fmt.Println("test schedule",time.Now())
		//time.Sleep( time.Duration(200)*time.Millisecond  )
		//t := time.Now()
		log.Println(time.Now().Format("20060102150405"))

		webdata.WebCmd(1)
		database.SetAcStatus2(2)
		/* reserved delay */
		StartTime := time.Now().Unix()
		DelayTime := interval - (StartTime % interval)
		//fmt.Println(DelayTime)
		time.Sleep(time.Duration(DelayTime) * time.Second)

	}
}

// Schedule2 check every second
func Schedule2(interval int64) {
	for {
		/* initial delay */
		//time.Sleep( time.Duration(1)*time.Second  )

		/* put the schedule code below */
		//fmt.Println("test schedule",time.Now())
		//time.Sleep( time.Duration(500)*time.Millisecond  )

		//fmt.Println("test schedule", time.Now())

		/* reserved delay */
		StartTime := time.Now().Unix()
		DelayTime := interval - (StartTime % interval)
		//fmt.Println(DelayTime)
		time.Sleep(time.Duration(DelayTime) * time.Second)
	}
}

// Schedule3 check every n second
func Schedule3(interval int64) {
	//for {
	/* initial delay */
	//time.Sleep( time.Duration(1)*time.Second  )

	/* put the schedule code here */
	//fmt.Println("test schedule",time.Now())
	//time.Sleep( time.Duration(500)*time.Millisecond  )

	log.Println("schedule3 top:", time.Now())

	/* reserved delay */
	StartTime := time.Now().Unix()
	DelayTime := interval - (StartTime % interval)
	//fmt.Println(DelayTime)
	time.Sleep(time.Duration(DelayTime) * time.Second)

	webdata.WebCmd(1)
	log.Println("schedule3 dwn:", time.Now())
	time.Sleep(time.Duration(2) * time.Second)

	//}
}

// Schedule4 check every 15 minutes
func Schedule4(interval int64, InitDelay int64) {
	/* initial delay */
	time.Sleep(time.Duration(InitDelay) * time.Second)
	for {

		/* put the schedule code here */
		//fmt.Println("test schedule",time.Now())
		//time.Sleep( time.Duration(200)*time.Millisecond  )
		//t := time.Now()
		log.Println(time.Now().Format("20060102150405"))

		webdata.WebCmd(2)
		database.SetAcStatus2(1)
		/* reserved delay */
		StartTime := time.Now().Unix()
		DelayTime := interval - (StartTime % interval)
		//fmt.Println(DelayTime)
		time.Sleep(time.Duration(DelayTime) * time.Second)

	}
}

// Schedule initial each schedule as thread
func Schedule() {
	/* schedule */
	StartTime := time.Now().Unix()

	log.Println(time.Now(), StartTime)

	interval := int64(60 * 15) //15 minutes interval
	DelayTime := interval - (StartTime % interval)

	log.Println("delay:", DelayTime)
	//time.Sleep( time.Duration(DelayTime)*time.Second  )

	go Schedule1(60*15, DelayTime) //every 15 minutes
	go Schedule2(1)                //schedule every 1 second routine

	go Schedule3(5) //one time schedule 1 second routine

}
