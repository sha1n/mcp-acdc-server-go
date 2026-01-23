---
name: code-review
description: Performs a detailed code review. Can review a specific commit or all local changes.
arguments:
  - name: commit
    description: The git commit hash to review.
    required: false
  - name: instructions
    description: Additional instructions or focus areas for the review.
    required: false
---
Please perform a detailed code review with the following context:

{{if .commit}}
1. Focus on the changes introduced in commit `{{.commit}}`.
{{else}}
1. Focus on all currently modified files in the repository.
{{end}}

{{if .instructions}}
**Additional Instructions:**
{{.instructions}}
{{end}}

Please provide feedback on code architecture, potential bugs, security concerns, and adherence to style guidelines.
