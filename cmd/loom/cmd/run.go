package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/orieken/loom/internal/orchestrator"
	"github.com/spf13/cobra"
)

type runFlags struct {
	spec          string
	resume        bool
	provider      string
	plan          string
	approve       string
	mockHangStage string
}

var runArgs runFlags

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute the delivery pipeline for a feature spec (experimental)",
	Long: `Execute the delivery pipeline for a feature spec with durable state.

EXPERIMENTAL SKELETON (roadmap M0.4, ADR-006): runs the built-in linear
deliver-feature plan stage by stage via the claude CLI, persisting
run-state.json after every transition. First Ctrl-C checkpoints and exits
cleanly (resume with --resume); a second Ctrl-C kills immediately.

Approval gates (roadmap L2.13) are process interrupts: the executor refuses
to start a gated stage until a human approves its gate. On a terminal you
are asked at the barrier; otherwise the run halts with exit code 3 and
prints the resume command (loom run --spec X --resume --approve <gate>).
Nothing an agent returns can approve a gate.

Not yet implemented (see BUILD-ROADMAP.md): reset-on-edit for approvals
(L2.14), retries/backoff, parallelism (L3.3), policy evaluation (L2.16),
conditional stage routing (L3.1), telemetry (L3.8).`,
	Args: cobra.NoArgs,
	RunE: runRun,
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVar(&runArgs.spec, "spec", "", "feature spec markdown file (required)")
	runCmd.Flags().BoolVar(&runArgs.resume, "resume", false, "continue an interrupted run from its checkpoint")
	runCmd.Flags().StringVar(&runArgs.approve, "approve", "", "approve the gate the run is waiting on (requires --resume)")
	runCmd.Flags().StringVar(&runArgs.provider, "provider", "claude", "stage provider: claude or mock")
	runCmd.Flags().StringVar(&runArgs.plan, "plan", orchestrator.DefaultDeliverFeaturePlanName, "pipeline plan to execute")
	runCmd.Flags().StringVar(&runArgs.mockHangStage, "mock-hang-stage", "", "mock provider only: stage ID that hangs until interrupted (testing)")
	_ = runCmd.Flags().MarkHidden("mock-hang-stage")
	_ = runCmd.MarkFlagRequired("spec")
}

func runRun(cmd *cobra.Command, _ []string) error {
	plan, err := selectPlan(runArgs.plan)
	if err != nil {
		return err
	}
	workspace, input, err := prepareRunWorkspace(runArgs.spec)
	if err != nil {
		return err
	}
	store := orchestrator.NewStateStore(filepath.Join(workspace, orchestrator.RunStateFileName))
	if err := checkResumeState(store, runArgs.resume); err != nil {
		return err
	}
	provider, err := selectProvider(plan, runArgs.provider, runArgs.mockHangStage)
	if err != nil {
		return err
	}
	return executeRun(cmd, plan, provider, store, input)
}

func executeRun(cmd *cobra.Command, plan orchestrator.Plan, provider orchestrator.Provider, store *orchestrator.StateStore, input orchestrator.StageInput) error {
	executor := orchestrator.NewExecutor(provider, store)
	executor.OnStale(func(stale []orchestrator.StaleStage) { reportStaleStages(cmd, stale) })
	executor.OnApprovalReset(func(reset *orchestrator.StaleApprovalError) { reportApprovalReset(cmd, reset) })
	if err := applyApproveFlag(executor, runArgs.approve, runArgs.resume); err != nil {
		return err
	}
	ctx, stopSignals := interruptibleContext(cmd)
	defer stopSignals()
	cmd.Printf("Running plan %q (%d stages) — state: %s\n", plan.Name, len(plan.Stages), store.Path())
	err := runWithGates(ctx, cmd, executor, plan, input)
	if errors.Is(err, context.Canceled) {
		cmd.Printf("Interrupted — checkpoint saved. Continue with: loom run --spec %s --resume\n", runArgs.spec)
		return err
	}
	if err != nil {
		return err
	}
	cmd.Printf("Plan %q completed. Artifacts: %s\n", plan.Name, input.WorkspaceDir)
	return nil
}

// interruptibleContext cancels on the first SIGINT/SIGTERM so the executor
// persists a clean checkpoint; a second signal kills the process immediately.
func interruptibleContext(cmd *cobra.Command) (context.Context, func()) {
	ctx, cancel := context.WithCancel(cmd.Context())
	signals := make(chan os.Signal, 2)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signals
		cmd.PrintErrln("interrupt received — checkpointing (press Ctrl-C again to kill)")
		cancel()
		<-signals
		os.Exit(130)
	}()
	return ctx, func() { signal.Stop(signals); cancel() }
}

func selectPlan(name string) (orchestrator.Plan, error) {
	if name != orchestrator.DefaultDeliverFeaturePlanName {
		return orchestrator.Plan{}, fmt.Errorf("unknown plan %q — only %q exists today (custom plans are a later roadmap item)",
			name, orchestrator.DefaultDeliverFeaturePlanName)
	}
	return orchestrator.DefaultDeliverFeaturePlan(), nil
}

// prepareRunWorkspace validates the spec and creates the feature workspace
// (.claude/feature-workspace/<feature>/, same location the markdown pipeline
// uses — see deliver-feature SKILL.md "Workspace Path Resolution").
func prepareRunWorkspace(specPath string) (string, orchestrator.StageInput, error) {
	absSpec, err := filepath.Abs(specPath)
	if err != nil {
		return "", orchestrator.StageInput{}, fmt.Errorf("resolve spec path: %w", err)
	}
	if _, err := os.Stat(absSpec); err != nil {
		return "", orchestrator.StageInput{}, fmt.Errorf("feature spec not found: %w", err)
	}
	feature := orchestrator.FeatureNameFromSpec(absSpec)
	workspace, err := filepath.Abs(filepath.Join(".claude", "feature-workspace", feature))
	if err != nil {
		return "", orchestrator.StageInput{}, fmt.Errorf("resolve workspace path: %w", err)
	}
	if err := os.MkdirAll(workspace, 0o755); err != nil {
		return "", orchestrator.StageInput{}, fmt.Errorf("create workspace: %w", err)
	}
	return workspace, orchestrator.StageInput{SpecPath: absSpec, WorkspaceDir: workspace}, nil
}

// checkResumeState enforces the resume contract: --resume requires existing
// state, and a fresh run refuses to start over existing state.
func checkResumeState(store *orchestrator.StateStore, resume bool) error {
	state, err := store.Load()
	if err != nil {
		return err
	}
	if resume && state == nil {
		return fmt.Errorf("--resume given but no run state exists at %s — start without --resume", store.Path())
	}
	if !resume && state != nil {
		return fmt.Errorf("run state already exists at %s — continue with --resume, or delete the file to start over", store.Path())
	}
	return nil
}
