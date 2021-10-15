package http

import (
	"net/http"

	"github.com/google/go-github/v39/github"
	"go.uber.org/zap"

	"github.com/tyuhara/team-review/config"
	gh "github.com/tyuhara/team-review/github"
)

func Handler(conf *config.Config, logger *zap.SugaredLogger) {
	http.HandleFunc("/github/events", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		payload, err := github.ValidatePayload(r, []byte(conf.GithubAppSecret))
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		event, err := github.ParseWebHook(github.WebHookType(r), payload)
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// ToDo: Move GetInstallation GetID and client creation section

		switch event := event.(type) {
		//// Triger when the new issue open or labels
		//case *github.IssuesEvent:
		//	switch event.GetAction() {
		//	case "labeled":
		//		logger.Infof("Get issue event: %v", event.GetAction())
		//		if err := gh.ProcessIssuesEvent(ctx, conf, event); err != nil {
		//			logger.Error(err)
		//			w.WriteHeader(http.StatusInternalServerError)
		//			return
		//		}
		//	case "opend":
		//		logger.Infof("Get issue event: %v", event.GetAction())
		//		if err := gh.ProcessIssuesEvent(ctx, conf, event); err != nil {
		//			logger.Error(err)
		//			w.WriteHeader(http.StatusInternalServerError)
		//			return
		//		}
		//	default:
		//		logger.Infof("Get issue event: %v", event.GetAction())
		//		if err := gh.ProcessIssuesEvent(ctx, conf, event); err != nil {
		//			logger.Error(err)
		//			w.WriteHeader(http.StatusInternalServerError)
		//			return
		//		}
		//	}
		// Trriger when the new comment add to the issue or pull request
		case *github.IssueCommentEvent:
			switch event.GetAction() {
			case "created":
				logger.Infof("Get issue comment event: %v", event.GetAction())
				if err := gh.HandleMergeIssueComment(ctx, event, conf, logger); err != nil {
					logger.Error(err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			default:
				logger.Infof("Get issue comment event: %v", event.GetAction())
			}
			//// Trriger when new pull request creates
			//case *github.PullRequestEvent:
			//	switch event.GetAction() {
			//	case "opened":
			//		logger.Infof("Get pull request event: %v", event.GetAction())
			//		if err := gh.ProcessPullRequestsEvent(ctx, conf, event); err != nil {
			//			logger.Error(err)
			//			w.WriteHeader(http.StatusInternalServerError)
			//			return
			//		}
			//	default:
			//		logger.Infof("Get pull request event: %v", event.GetAction())
			//	}
			//default:
			//	logger.Info("Get request")
		}
	})

	logger.Info("[INFO] Server listening")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Error(err)
	}
}
