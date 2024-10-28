Output images/component-defn-console-validate.gif

Require lula
Set FontSize 14
Set Width 1850
Set Height 925
Set Framerate 24
Set Padding 5

Hide
Type "lula console -f ./src/test/unit/common/oscal/valid-multi-component-validations.yaml -a tmp.yaml" Enter
Sleep 1s
Show

# Select the framework
Sleep 500ms
Right
Sleep 500ms
Right
Sleep 500ms
Enter
Sleep 500ms
Down
Sleep 500ms
Down
Sleep 500ms
Enter
Sleep 1s

# Run Validate
Ctrl+v
Sleep 2s
Enter
Sleep 6s

Hide
Ctrl+c
Sleep 50ms
Enter
Sleep 50ms
Type "rm tmp.yaml" Enter