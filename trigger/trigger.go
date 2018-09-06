package trigger

import (
	"fmt"
	"path"
	"path/filepath"
	"regexp"

	"github.com/NeoJRotary/GCB-bridge/app"
	D "github.com/NeoJRotary/describe-go"
)

// DEBUG enable DEBUG to print log of trigger
var DEBUG = D.GetENV("DEBUG", "") == "true"

type trigger struct {
	cb *cloudBuild

	Branches         []string `yaml:"branches"`
	Tags             []string `yaml:"tags"`
	PullRequestBases []string `yaml:"pullRequestBases"`
	IncludedFiles    []string `yaml:"includedFiles"`
	IgnoredFiles     []string `yaml:"ignoredFiles"`
}

func (trg *trigger) debugLog(a ...interface{}) {
	if !DEBUG {
		return
	}
	trg.cb.Log.Push(fmt.Sprint(a...))
}

func (trg *trigger) isValid(cb *cloudBuild) bool {
	trg.cb = cb

	// if there is no any main trigger (branch/tag/PR) then just check files
	if len(trg.Branches)+len(trg.Tags)+len(trg.PullRequestBases) == 0 {
		trg.debugLog("isValid: no any main trigger")
		return trg.validFiles(cb.repo)
	}

	// check by Event
	switch cb.repo.Event {
	case "Branch":
		if !trg.validBranch(cb.repo.Branch) && !trg.validAssociatedBases(cb.repo.AssociatedBases) {
			return false
		}
	case "Tag":
		if !trg.validTag(cb.repo.Tag) {
			return false
		}
	case "PullRequest":
		if !trg.validPullRequest(cb.repo.BaseBranch) {
			return false
		}
	}

	return trg.validFiles(cb.repo)
}

func (trg *trigger) validBranch(branch string) bool {
	// doesn't match basic event requirement
	if len(trg.Branches) == 0 || branch == "" {
		trg.debugLog("validBranch: doesn't match basic event requirement")
		return false
	}

	trg.debugLog("validBranch: branch[", branch, "] trg.Branches:", trg.Branches)

	// go through all regex
	for _, reg := range trg.Branches {
		if regexp.MustCompile(reg).MatchString(branch) {
			return true
		}
	}

	trg.debugLog("validBranch: failed")

	return false
}

func (trg *trigger) validAssociatedBases(bases []string) bool {
	// doesn't match basic event requirement
	if len(trg.PullRequestBases) == 0 || len(bases) == 0 {
		trg.debugLog("validAssociatedBases: doesn't match basic event requirement")
		return false
	}

	trg.debugLog("validAssociatedBases: bases:", bases, "PullRequestBases:", trg.PullRequestBases)

	// go through all regex
	for _, reg := range trg.PullRequestBases {
		for _, base := range bases {
			if regexp.MustCompile(reg).MatchString(base) {
				return true
			}
		}
	}

	trg.debugLog("validAssociatedBases: failed")

	return false
}

func (trg *trigger) validTag(tag string) bool {
	// doesn't match basic event requirement
	if len(trg.Tags) == 0 || tag == "" {
		trg.debugLog("validTag: doesn't match basic event requirement")
		return false
	}

	trg.debugLog("validTag: tag[", tag, "] trg.Tags:", trg.Tags)

	// go through all regex
	for _, reg := range trg.Tags {
		if regexp.MustCompile(reg).MatchString(tag) {
			return true
		}
	}

	trg.debugLog("validTag: failed")

	return false
}

func (trg *trigger) validPullRequest(base string) bool {
	// doesn't match basic event requirement
	if len(trg.PullRequestBases) == 0 || base == "" {
		trg.debugLog("validPullRequest: doesn't match basic event requirement")
		return false
	}

	trg.debugLog("validPullRequest: base[", base, "] trg.PullRequestBases:", trg.PullRequestBases)

	// go through all regex
	for _, reg := range trg.PullRequestBases {
		if regexp.MustCompile(reg).MatchString(base) {
			return true
		}
	}

	trg.debugLog("validPullRequest: failed")

	return false
}

func (trg *trigger) validFiles(repo *app.Repo) bool {
	// if no files conf, return true
	if len(trg.IncludedFiles) == 0 {
		trg.debugLog("validFiles: IncludedFiles is empty, pass")
		return true
	}

	changes := repo.GetChanges()
	// cannot get chganes, always return false
	if changes == nil {
		trg.debugLog("validFiles: failed, no any changes")
		return false
	}

	trg.debugLog("validFiles: changes", changes)
	changesD := D.StringSlice(changes)

	for _, ignore := range trg.IgnoredFiles {
		matches, err := filepath.Glob(path.Join(repo.Dir, ignore))
		if D.IsErr(err) {
			trg.debugLog("validFiles: matching", err)
			continue
		}

		trg.debugLog("validFiles: IGNORE > ", ignore, matches)

		changesD.DeleteSame(matches...)
	}

	// all changes be ignored
	if changesD.Empty() {
		trg.debugLog("validFiles: failed because all changes be ignored")
		return false
	}

	for _, include := range trg.IncludedFiles {
		matches, err := filepath.Glob(path.Join(repo.Dir, include))
		if D.IsErr(err) {
			trg.debugLog("validFiles: matching", err)
			continue
		}

		trg.debugLog("validFiles: INCLUDE > ", include, matches)

		// find includes
		if changesD.Include(matches...) {
			return true
		}
	}

	trg.debugLog("validFiles: nothing matched")

	// nothing match
	return false
}
