Output images/assessment-results-console.gif

Require lula
Set FontSize 14
Set Width 1850
Set Height 925
Set Framerate 24
Set Padding 5

Hide
Type "lula console -f ./src/test/unit/common/oscal/valid-assessment-results-multi.yaml" Enter
Sleep 1s
Tab
Sleep 500ms
Show

# Show results picker
Sleep 1s
Right
Sleep 500ms
Enter
Sleep 1s
Down
Sleep 500ms
Enter
Sleep 500ms

# Navigate to findings, filter
Right
Sleep 500ms
Right
Sleep 500ms
Down
Sleep 500ms
Up
Sleep 500ms
Type "/"
Sleep 500ms
Type "ID-1"
Sleep 500ms
Enter
Sleep 1s

# Navigate to observations
Right
Sleep 500ms
Type "d"
Sleep 2s

