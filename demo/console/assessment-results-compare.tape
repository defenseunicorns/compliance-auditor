Output images/assessment-results-console-compare.gif

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

# Select comparison result
Sleep 1s
Right
Sleep 500ms
Right
Sleep 500ms
Enter
Sleep 1s
Down
Sleep 500ms
Enter
Sleep 500ms

# Navigate to findings, observations
Right
Sleep 500ms
Right
Sleep 500ms
Down
Sleep 500ms
Up
Sleep 500ms
Left
Sleep 1s

# Show findings comparison detail
Type "d"
Sleep 2s

