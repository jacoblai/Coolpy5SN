package Dispatcher

import (
	"fmt"
	"io/ioutil"
	"bytes"
	"net/http"
	"strconv"
)

// NewWorker creates takes a numeric id and a channel w/ worker pool.
func NewWorker(id int, cphttp string, workerPool chan chan Job) Worker {
	return Worker{
		id:         id,
		jobQueue:   make(chan Job),
		workerPool: workerPool,
		quitChan:   make(chan bool),
		cpHttp:cphttp,
	}
}

type Worker struct {
	id         int
	jobQueue   chan Job
	workerPool chan chan Job
	quitChan   chan bool
	cpHttp     string
}

func (w Worker) Start() {
	go func() {
		for {
			// Add my jobQueue to the worker pool.
			w.workerPool <- w.jobQueue

			select {
			case job := <-w.jobQueue:
				url :=w.cpHttp + "/api/hub/" + strconv.FormatInt(job.HubId, 10) + "/node/" + strconv.FormatInt(job.NodeId, 10) + "/datapoints"
				fmt.Println("URL:>", url)

				req, _ := http.NewRequest("POST", url, bytes.NewBuffer(job.CpJson))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("U-ApiKey", job.Ukey)

				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					fmt.Println(err)
				}
				defer resp.Body.Close()

				fmt.Println("response Status:", resp.Status)
				body, _ := ioutil.ReadAll(resp.Body)
				fmt.Println("response Body:", string(body))
			case <-w.quitChan:
				fmt.Printf("worker%d stopping\n", w.id)
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.quitChan <- true
	}()
}
