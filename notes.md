# Profile Support

Processes that include profile

- generate a comp-def using a profile as source
- comp-def control-implementation source
- generate a catalog (profile resolve?) 
- generate a profile

## Generate Component Support

If the source of an intended generation is a `profile`, we will need to:
- Create a profile object
- Attempt to identify a catalog
  - If found - resolve catalog first -> then perform the mapping
  - If not found?

## Generate a Profile

Source Location
- We can inject this directly in the `href` for creation purposes and opinionation
- When processing the import `href` we will need to account for back-matter reference
  - More so the OSCAL content examples point to a back-matter item with multiple links to each `mediatype` option of xml/json/yaml
  - Consider processing mediatype

Include OR Exclude controls
- Likely will not perform an `AND` operation here. 

Merge
- Expose as a flag
- initially support `as-is` only
  - flat is an `empty` type which is hard to represent
  - custom introduces a lot of complexity and difficult to find good examples
  - combine includes a vague description and lacks example of when conflicting ID's would occur

## Profile Resolve

