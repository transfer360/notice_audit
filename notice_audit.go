package notice_audit

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	joonix "github.com/joonix/log"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

type Package struct {
	Sref             string    `json:"sref,omitempty"`
	SoftwareProvider int       `json:"software_provider,omitempty"`
	FleetSoftware    int       `json:"fleet_software,omitempty"`
	DateTime         time.Time `json:"date_time"`
	Action           string    `json:"action"`
	SystemAction     bool      `json:"system_action"`
}

func (p *Package) Send(ctx context.Context) error {

	if len(os.Getenv("DEVELOPMENT")) == 0 {
		log.SetFormatter(joonix.NewFormatter())
	}

	log.SetLevel(log.DebugLevel)

	client, err := pubsub.NewClient(ctx, "transfer-360")
	if err != nil {
		log.Fatal(err)
	}

	jsonData, err := json.Marshal(p)
	if err != nil {
		log.Errorln(err)
		return err
	}

	topic := client.Topic("notice_audit_actions")

	msg := &pubsub.Message{
		Data: jsonData,
	}

	if _, err := topic.Publish(ctx, msg).Get(ctx); err != nil {
		log.Errorln(err)
		return err
	}

	return nil

}
