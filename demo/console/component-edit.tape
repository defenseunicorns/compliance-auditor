Output images/component-defn-console-edit.gif

Require lula
Set FontSize 14
Set Width 1850
Set Height 925
Set Framerate 24
Set Padding 5

Hide
Type "lula console -f ./src/test/unit/common/oscal/valid-multi-component-validations.yaml" Enter
Sleep 1s
Show

# Navigate and select a control
Right
Sleep 500ms
Right
Sleep 500ms
Right
Sleep 1s
Enter
Sleep 1s

# Navigate to description
Right
Sleep 500ms
Right
Sleep 500ms
Type "e"
Sleep 1s

# Add new line and write some text
Ctrl+E
Sleep 500ms
Type "Here is some text"
Sleep 500ms
Enter
Sleep 2s

# Show save dialog
Ctrl+S
Sleep 2s
