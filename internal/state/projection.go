package state

// A projection is the narrowly-scoped slice of upstream state a stage is
// allowed to read — the point of L2.9. The architect does not receive the
// whole analysis; it receives the fields architecture-contract.md needs to
// make structural decisions.
//
// Applying projections across every stage, and retiring the LLM
// summarization currently on the inter-stage path, is L2.10.

import (
	"encoding/json"
	"fmt"
)

// ArchitectInput is what the architect reads from the analysis: the
// constraints and structural facts it decides against. Notably absent are
// the QA, tech-writer, and devops task lists, the definition of done, and
// the edge cases — the architect has no use for them, so it never sees
// them.
type ArchitectInput struct {
	Feature                   string                     `json:"feature"`
	Summary                   string                     `json:"summary"`
	AcceptanceCriteria        []AcceptanceCriterion      `json:"acceptanceCriteria"`
	NonFunctionalRequirements []NonFunctionalRequirement `json:"nonFunctionalRequirements,omitempty"`
	ProposedFitnessFunctions  []FitnessFunction          `json:"proposedFitnessFunctions,omitempty"`
	BoundedContext            BoundedContext             `json:"boundedContext"`
	DomainEvents              DomainEvents               `json:"domainEvents,omitempty"`
	AffectedComponents        []AffectedComponent        `json:"affectedComponents"`
	DataModelChanges          []DataModelChange          `json:"dataModelChanges,omitempty"`
	APIChanges                []APIChange                `json:"apiChanges,omitempty"`
	NewDependencies           []string                   `json:"newDependencies,omitempty"`
	ArchitecturalFlags        []string                   `json:"architecturalFlags,omitempty"`
	DeveloperTasks            []string                   `json:"developerTasks,omitempty"`
	OpenQuestions             []string                   `json:"openQuestions,omitempty"`
}

// DeveloperInput is what the developer reads when the loop sends it round
// again: the verdict and the findings it has to address. Notably absent is
// everything the reviewer said about surfaces and scores — the developer's
// next iteration acts on instructions, not on an assessment of its last
// attempt.
type DeveloperInput struct {
	Feature  string    `json:"feature"`
	Verdict  Verdict   `json:"verdict"`
	Findings []Finding `json:"findings"`
}

// ProjectReviewForDeveloper narrows a review to what the next iteration
// needs. It is a field selection, like every other projection: no model
// summarises the findings on the way back.
func ProjectReviewForDeveloper(review ReviewState) DeveloperInput {
	return DeveloperInput{Feature: review.Feature, Verdict: review.Verdict, Findings: review.Findings}
}

// ProjectFor computes the projection a consumer of the given kind receives
// from its upstream document. It is a field selection, never a
// summarization: no model sits on this path.
func ProjectFor(consumer Kind, upstream Kind, upstreamPayload []byte) ([]byte, error) {
	if upstream == KindReview {
		return projectReview(upstreamPayload)
	}
	if consumer != KindArchitecture || upstream != KindAnalysis {
		return nil, fmt.Errorf("no projection declared from %q to %q", upstream, consumer)
	}
	decoded, err := Decode(KindAnalysis, upstreamPayload)
	if err != nil {
		return nil, fmt.Errorf("upstream analysis state: %w", err)
	}
	analysis, ok := decoded.(*AnalysisState)
	if !ok {
		return nil, fmt.Errorf("upstream state is not an analysis")
	}
	return json.Marshal(projectAnalysisForArchitect(*analysis))
}

func projectReview(payload []byte) ([]byte, error) {
	decoded, err := Decode(KindReview, payload)
	if err != nil {
		return nil, fmt.Errorf("upstream review state: %w", err)
	}
	review, ok := decoded.(*ReviewState)
	if !ok {
		return nil, fmt.Errorf("upstream state is not a review")
	}
	return json.Marshal(ProjectReviewForDeveloper(*review))
}

func projectAnalysisForArchitect(analysis AnalysisState) ArchitectInput {
	return ArchitectInput{
		Feature: analysis.Feature, Summary: analysis.Summary,
		AcceptanceCriteria:        analysis.AcceptanceCriteria,
		NonFunctionalRequirements: analysis.NonFunctionalRequirements,
		ProposedFitnessFunctions:  analysis.ProposedFitnessFunctions,
		BoundedContext:            analysis.BoundedContext,
		DomainEvents:              analysis.DomainEvents,
		AffectedComponents:        analysis.AffectedComponents,
		DataModelChanges:          analysis.DataModelChanges,
		APIChanges:                analysis.APIChanges,
		NewDependencies:           analysis.NewDependencies,
		ArchitecturalFlags:        analysis.ArchitecturalFlags,
		DeveloperTasks:            analysis.Tasks.Developer,
		OpenQuestions:             analysis.OpenQuestions,
	}
}
