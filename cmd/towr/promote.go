package main

import (
	"fmt"
	"time"

	"github.com/brianaffirm/towr/internal/store"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

func newPromoteCmd(initApp func() (*appContext, error), jsonFlag *bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "promote <workspace-id>",
		Short:             "Promote to a workspace for human interaction",
		Long:              "Emit a task.promoted event and attach to the workspace's tmux session for direct human interaction.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: workspaceIDCompletion(initApp),
		RunE: func(cmd *cobra.Command, args []string) error {
			wsID := args[0]

			app, err := initApp()
			if err != nil {
				return err
			}

			// 1. Check for active dispatch — if running, refuse.
			latestDisp, err := app.store.LatestDispatch(app.repoRoot, wsID)
			if err != nil {
				return fmt.Errorf("check latest dispatch: %w", err)
			}
			if latestDisp != nil {
				dispID, _ := latestDisp.Data["dispatch_id"].(string)
				if dispID != "" {
					latestEvt, err := app.store.LatestTaskEvent(app.repoRoot, wsID, dispID)
					if err != nil {
						return fmt.Errorf("check latest task event: %w", err)
					}
					if latestEvt != nil && latestEvt.Kind != store.EventTaskCompleted && latestEvt.Kind != store.EventTaskFailed {
						return fmt.Errorf("workspace %q has active dispatch %s — wait for it to complete or report before promoting", wsID, dispID)
					}
				}
			}

			// 2. Emit task.promoted event.
			if err := app.store.EmitEvent(store.Event{
				ID:          uuid.New().String(),
				Kind:        store.EventTaskPromoted,
				WorkspaceID: wsID,
				RepoRoot:    app.repoRoot,
				Timestamp:   time.Now().UTC(),
				Data: map[string]interface{}{
					"actor": "human",
				},
			}); err != nil {
				return fmt.Errorf("emit promote event: %w", err)
			}

			// 3. Attach to tmux session.
			if err := app.term.Attach(wsID); err != nil {
				return fmt.Errorf("attach to workspace: %w", err)
			}

			return nil
		},
	}

	return cmd
}
