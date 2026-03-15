package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/brianaffirm/towr/internal/control"
	"github.com/brianaffirm/towr/internal/orchestrate"
	"github.com/spf13/cobra"
)

func newRunCmd(initApp func() (*appContext, error), jsonFlag *bool) *cobra.Command {
	var budgetOverride float64
	var quiet, dryRun bool

	cmd := &cobra.Command{
		Use:   "run <plan.yaml>",
		Short: "Execute a plan: spawn, dispatch, approve, PR, watch — all in one",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			app, err := initApp()
			if err != nil {
				return err
			}
			plan, err := orchestrate.LoadPlan(args[0])
			if err != nil {
				return err
			}
			if err := plan.Validate(); err != nil {
				return fmt.Errorf("invalid plan: %w", err)
			}
			if plan.Settings.LandPR && !plan.Settings.CreatePR {
				plan.Settings.CreatePR = true
			}
			if budgetOverride > 0 {
				plan.Settings.Budget = budgetOverride
			}

			svc := &control.RunService{Store: app.store, Runtime: &controlRuntime{app: app},
				Router: &control.RouterAdapter{}, Clock: time.Now, Logger: &stdLog{}}
			req := buildRunRequest(app.repoRoot, plan)

			if dryRun {
				fmt.Print(formatDryRun(plan.Name, svc.DryRun(req)))
				return nil
			}
			if !quiet {
				fmt.Print(formatDryRun(plan.Name, svc.DryRun(req)))
				if !plan.Settings.AutoApprove {
					fmt.Print("\nProceed? [Y/n] ")
					var answer string
					fmt.Scanln(&answer)
					if answer != "" && strings.ToLower(answer) != "y" {
						fmt.Println("Aborted.")
						return nil
					}
				}
			}

			_ = svc.ReconcileAll(context.Background(), app.repoRoot)
			if plan.Settings.Web {
				startWebDashboard(plan.Settings.WebAddr)
			}
			startMuxStatusUpdater()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			go func() {
				sigCh := make(chan os.Signal, 1)
				signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
				<-sigCh
				cancel()
			}()

			handle, err := svc.Start(ctx, req)
			if handle != nil {
				for handle.Status == control.RunRunning {
					time.Sleep(100 * time.Millisecond)
				}
				fmt.Printf("\nRun %s: %s\n", handle.ID, handle.Status)
			}
			if err == nil && plan.Settings.ReactToReviews {
				runWatchReact(app, parsePollInterval(plan))
			}
			return err
		},
	}

	cmd.Flags().Float64Var(&budgetOverride, "budget", 0, "Maximum USD budget (0 = no limit)")
	cmd.Flags().BoolVar(&quiet, "quiet", false, "Skip pre-run routing summary")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show routing table, then exit")
	return cmd
}
