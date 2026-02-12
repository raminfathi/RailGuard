package services

import (
	"fmt"
	"railguard/internal/core/domain"
	"railguard/internal/core/ports"
)

type SafetyValidatorService struct {
	rulesMap map[string]map[string]string // Cache rules in memory: map[CodeA][CodeB] -> Status
}

func NewSafetyValidatorService(repo ports.RuleRepository) (*SafetyValidatorService, error) {
	// Load all rules into memory at startup
	rulesList, err := repo.GetAllDangerRules()
	if err != nil {
		return nil, err
	}

	// Transform list into a fast lookup map
	rulesMap := make(map[string]map[string]string)
	for _, r := range rulesList {
		if rulesMap[r.CodeA] == nil {
			rulesMap[r.CodeA] = make(map[string]string)
		}
		rulesMap[r.CodeA][r.CodeB] = r.Status
	}

	return &SafetyValidatorService{rulesMap: rulesMap}, nil
}

// ValidateComposition checks the train order against dangerous goods matrix
func (v *SafetyValidatorService) ValidateComposition(wagons []domain.SelectedWagon) (bool, string) {

	// Iterate through all wagons to find pairs of dangerous goods
	for i := 0; i < len(wagons); i++ {
		// If wagon i has no dangerous goods, skip
		if !wagons[i].HasDangerousGoods {
			continue
		}

		// Check against all subsequent wagons
		for j := i + 1; j < len(wagons); j++ {
			// If wagon j has no dangerous goods, it acts as a buffer, just continue searching
			if !wagons[j].HasDangerousGoods {
				continue
			}

			// Both wagon i and j have dangerous goods. Check compatibility.
			codeA := wagons[i].DangerousGoodsCode
			codeB := wagons[j].DangerousGoodsCode

			status := v.getRuleStatus(codeA, codeB)
			distance := j - i // Difference in index (1 means adjacent)

			// Logic based on your matrix definition
			switch status {
			case "-":
				// Cannot be adjacent
				if distance == 1 {
					return false, fmt.Sprintf("Conflict: Wagon #%d (%s) cannot be adjacent to Wagon #%d (%s)",
						wagons[i].WagonSpec.Number, codeA, wagons[j].WagonSpec.Number, codeB)
				}
			case "1":
				// Needs 1 buffer wagon.
				// Indices: [i, i+1] -> diff=1 (Fail). [i, buffer, j] -> diff=2 (Pass)
				if distance < 2 {
					return false, fmt.Sprintf("Conflict: Wagon #%d (%s) needs 1 buffer wagon from Wagon #%d (%s)",
						wagons[i].WagonSpec.Number, codeA, wagons[j].WagonSpec.Number, codeB)
				}
			case "2":
				// Needs 2 buffer wagons.
				// Indices: [i, buff, j] -> diff=2 (Fail). [i, buff, buff, j] -> diff=3 (Pass)
				if distance < 3 {
					return false, fmt.Sprintf("Conflict: Wagon #%d (%s) needs 2 buffer wagons from Wagon #%d (%s)",
						wagons[i].WagonSpec.Number, codeA, wagons[j].WagonSpec.Number, codeB)
				}
			case "*", "+":
				// Allowed
				continue
			}
		}
	}

	return true, "All checks passed."
}

func (v *SafetyValidatorService) getRuleStatus(a, b string) string {
	if inner, ok := v.rulesMap[a]; ok {
		if status, ok := inner[b]; ok {
			return status
		}
	}
	// Fallback or assume forbidden if not found? usually "-"
	return "-"
}
