package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/go-co-op/gocron"

	"github.com/howarddo2208/cronman/config"
	"github.com/howarddo2208/cronman/models"
)

func scheduleJob(s *gocron.Scheduler, name string, job models.Job, stdout io.Writer) error {
	log.SetOutput(stdout)

	_, err := s.CronWithSeconds(job.Schedule).Do(func() {
		cmd := exec.Command("sh", "-c", job.Cmd)
		out, errCmd := cmd.Output()
		if errCmd != nil {
			fmt.Fprintf(os.Stderr, "%s\n", errCmd)
		} else {
			log.Printf("output of job %v: %s\n", name, out)
		}
	})
	if err != nil { // error scheduling
		return err
	}
	return nil
}

func runCrons(c *models.Config, stdout io.Writer) error {
	s := gocron.NewScheduler(time.Local)

	for name, job := range c.Jobs {
		if errSchedule := scheduleJob(s, name, job, stdout); errSchedule != nil {
			return errSchedule
		}
	}

	s.StartBlocking()
	return nil
}

func main() {
	c := &models.Config{}

	var err error
	if err = config.InitConfig(c); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}

	if err = runCrons(c, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
}
